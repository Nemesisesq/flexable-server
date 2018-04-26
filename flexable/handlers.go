package flexable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/globalsign/mgo/bson"
	"github.com/manveru/faker"
	"github.com/nemesisesq/flexable/account"
	PlivoClient "github.com/nemesisesq/flexable/plivio/client"
	"github.com/nemesisesq/flexable/position"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/nemesisesq/flexable/utils"
	"github.com/odknt/go-socket.io"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func OpenShiftHandler(s socketio.Conn, _ interface{}) {
	log.Info("Returning openshifts")

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User);
	var query bson.M
	//var user account.User
	//if !ok {
	//	fmt.Println("something is not ok")
	//}
	
	if user.Profile.Company.UUID == "" {
		user.Profile.Company.UUID = "123"
	}
	query = bson.M{"company.uuid": user.Profile.Company.UUID}
	shift_list := shifts.GetAllShifts(query)

	s.Emit(constructSocketID(OPEN_SHIFTS), shift_list, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
}

const NEW_SHIFT_TITLE = "There's a new shift!!!!"

func FindShiftReplacementHandler(s socketio.Conn, data interface{}) {
	log.Info("Finding a shift replacement")
	payload := data.(map[string]interface{})["payload"]

	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User)

	var shift shifts.Shift
	err = json.Unmarshal(tmp, &shift)

	if err != nil {
		panic(err)
	}

	shift.SmsID = uuid.NewV4().String()

	shift.Company = user.Profile.Company

	plivoClient, err := PlivoClient.NewClient()

	if err != nil {

		panic(err)
	}

	app := plivoClient.CreateApplication(shift.Company.Name, shift.SmsID)

	shift.Application = app
	//TODO
	//shift.Save()
	shift.Name = fmt.Sprintf("%v : %v per hour", shift.Job.Title, shift.Job.Compensation)
	// TODO uncomment for testing buying phone number works
	//number, err := plivio.BuyPhoneNumber(shift)

	response, err := plivoClient.Numbers.Update(
		"16143636301",
		plivo.NumberUpdateParams{
			AppID: app.AppID,
		},
	)

	log.WithFields(log.Fields{"response": response, "number": "some number to be determined", "app_id ": app.AppID}).Info("number updated to appID")

	if err != nil {
		panic(err)
	}

	//shift.PhoneNumber = number
	shift.PhoneNumber = "16143636301"

	shift.Save()

	templateString := `
Hey There is an open shift from {{.StartTime }} to {{.EndTime}}
On {{.Date }}. Reply "1" if you would like to pick up this shift. Skip the text messages and download the Flexable app in the Apple App store
or the Google play store.
`
	// The New Implementation
	users, err := account.FindAll(bson.M{"user.profile.company.uuid": shift.Company.UUID})
	buf, err := CreateTextMessageString(templateString, shift)
	for _, user := range users {
		user.Notify(buf.String(), NEW_SHIFT_TITLE, shift.PhoneNumber)
	}

	// The old implementation
	plivoClient.SendMessages(shift.PhoneNumber, "12165346715<16142881847", buf.String())
	if err != nil {
		panic(err)
	}
}
func CreateTextMessageString(templateString string, shift shifts.Shift) (bytes.Buffer, error) {
	messageTemplate, err := template.New("test").Parse(templateString)
	if err != nil {
		panic(err)
	}
	buf := bytes.Buffer{}
	err = messageTemplate.Execute(&buf, shift)
	if err != nil {
		panic(err)
	}
	return buf, err
}

type Data struct {
	Payload struct {
		Shift struct {
			V  int    `json:"__v"`
			ID string `json:"_id"`
			AbsentWorker struct {
				Name string `json:"name"`
			} `json:"absentWorker"`
			Application struct {
				APIID   string `json:"api_id"`
				AppID   string `json:"app_id"`
				Message string `json:"message"`
			} `json:"application"`
			Company struct {
				Name string `json:"name"`
				UUID string `json:"uuid"`
			} `json:"company"`
			Date string `json:"date"`
			End  string `json:"end"`
			Job struct {
				Compensation int    `json:"compensation"`
				Title        string `json:"title"`
			} `json:"job"`
			Location    string `json:"location"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			SmsID       string `json:"sms_id"`
			Start       string `json:"start"`
			Volunteers struct {
				Number string `json:"number"`
			} `json:"volunteers"`
		} `json:"shift"`
		Volunteer account.User `json:"volunteer"`
	} `json:"payload"`
}

func SelectVolunteer(s socketio.Conn, data interface{}) interface{} {
	log.Info("Selecting the Volunteer")
	//spew.Dump(data)

	payload := &Data{}
	tmp, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(tmp, &payload)

	p := payload.Payload
	shift := shifts.GetOneShift(bson.M{"sms_id": p.Shift.SmsID})

	shift.Chosen = payload.Payload.Volunteer

	plivoClient, err := PlivoClient.NewClient()

	if err != nil {
		panic(err)
	}

	t := `
				Awesome!!! You've picked up a shift!
				Details:

				Location : {{.Location}}
				Date: {{.Date}}
				Start Time : {{.StartTime}}
				End Time : {{.EndTime }}
				`

	buf, err := CreateTextMessageString(t, shift)

	if err != nil {
		panic(err)
	}

	err = plivoClient.SendMessages(shift.PhoneNumber, shift.Chosen.Profile.PhoneNumber, buf.String())
	if err != nil {
		panic(err)
	}

	shift.Save()

	return nil
}

func GetAvailableEmployees(s socketio.Conn, data interface{})  {

	log.Info("Getting Available Employees")
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	empList, err := account.FindAll(nil)

	if err != nil {
		panic(err)
	}
	if len(empList) <= 1 {

		for i := 0; i < utils.RandomRange(2, 10); i++ {
			//var num string
			//bin := utils.RandomRange(1, 2)

			//if bin%2 == 0 {
			//	num = "12165346715"
			//} else {
			//	num = "16142881847"
			//}

			x := account.User{
				ID: bson.NewObjectId(),
				//Name:   fake.Name(),
				//Number: num,
				Email: fake.Email(),
				Profile: account.Profile{Location: account.GeoLocation{
					fake.Latitude(),
					fake.Longitude(),
				},
				},
				Position: position.Position{
					ID:           bson.NewObjectId(),
					Title:        fake.JobTitle(),
					Compensation: 10.00,
					Rate:         "hr",
				},
			}

			x.Upsert(bson.M{"_id": x.ID})
			empList = append(empList, x)

		}
	}

	fmt.Println("sending employees")

	log.Info(s.Namespace())

	id := constructSocketID(EMPLOYEE_LIST)
	s.Emit(id, empList, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
}

func GetPositions(s socketio.Conn, data interface{}) {
	log.Info("Getting positions")
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	jobs := position.GetAllPositions(nil)

	if len(jobs) <= 1 {

		for i := 0; i < utils.RandomRange(2, 10); i++ {
			x := position.Position{
				bson.NewObjectId(),
				fake.JobTitle(),
				10.00,
				"hr",
			}

			x.Save()
			jobs = append(jobs, x)
		}
	}

	fmt.Println(jobs)

	s.Emit(constructSocketID(GET_JOBS), jobs, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})

}
