package flexable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/globalsign/mgo/bson"
	"github.com/manveru/faker"
	"github.com/nemesisesq/flexable-server/account"
	PlivoClient "github.com/nemesisesq/flexable-server/plivio/client"
	"github.com/nemesisesq/flexable-server/position"
	"github.com/nemesisesq/flexable-server/shifts"
	"github.com/nemesisesq/flexable-server/utils"
	"github.com/odknt/go-socket.io"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/mitchellh/hashstructure"
	"github.com/jinzhu/now"
	"github.com/nemesisesq/flexable-server/messaging"
)

var st time.Time

func OpenShiftHandler(s socketio.Conn, _ interface{}) {
	st = time.Now()
	log.Info("Returning openshifts")

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User);
	var query bson.M
	if user.Profile.Company.UUID == "" {
		user.Profile.Company.UUID = "123"
	}
	query = bson.M{"company.uuid": user.Profile.Company.UUID}

	ticker := time.NewTicker(time.Second * 2)
	timeout := time.NewTimer(time.Minute)
	var currentShiftState uint64
	go func() {
		shiftList := []shifts.Shift{}

		timeout, currentShiftState = emitCurrentShifts(shiftList, query, currentShiftState, s, timeout)
	L:

		for {

			user = *user.Find(bson.M{"_id": user.ID})
			shiftList := []shifts.Shift{}
			select {
			case <-ticker.C:
				timeout, currentShiftState = emitCurrentShifts(shiftList, query, currentShiftState, s, timeout)

			case <-ctx.Done():
				ticker.Stop()
				break L

			case <-timeout.C:
				s.Close()
				cancel := ctx.Value("cancel").(context.CancelFunc)
				cancel()
			}
		}
	}()
}

func emitCurrentShifts(shiftList []shifts.Shift, query bson.M, currentShiftState uint64, s socketio.Conn, timeout *time.Timer) (*time.Timer, uint64) {
	shiftList = shifts.GetAllShifts(query)
	cleaned_shift_list := []shifts.SkinnyShift{}
	for _, v := range shiftList {
		present := time.Now().AddDate(0, 0, -7)
		date := now.MustParse(v.Date)

		if present.Before(date) {
			x := &shifts.SkinnyShift{}

			x.ID = v.ID
			x.Name = v.Name
			x.Job = v.Job
			x.Date = v.Date
			x.StartTime = v.StartTime
			x.EndTime = v.EndTime
			x.Volunteers = len(v.Volunteers)
			if v.Chosen.Profile.Email != "" {

				x.Chosen = true
			}
			cleaned_shift_list = append(cleaned_shift_list, *x)
		}
	}
	shift_list_hash, err := hashstructure.Hash(&cleaned_shift_list, nil)
	if err != nil {
		panic(err)
	}
	if currentShiftState != shift_list_hash {
		log.Info(currentShiftState, shift_list_hash)
		currentShiftState = shift_list_hash

		s.Emit(constructSocketID(OPEN_SHIFTS), cleaned_shift_list)
		finished := time.Now()

		elapsed := st.Sub(finished)

		log.WithField("elapsed_time ", elapsed.Seconds()).Info("open shift")
		timeout = time.NewTimer(time.Minute)
	}
	return timeout, currentShiftState
}

const NEW_SHIFT_TITLE = "There's a new shift!!!!"

func GetCompanyList(s socketio.Conn, data interface{}){
	log.Info("Getting Company List")
	rpc := messaging.NewRpcCient(messaging.COMPANY_RPI_QUEUE, func(message *messaging.RpcMessage) interface{} {
		log.Info("Returning Payload ")
		s.Emit(constructSocketID(GET_COMPANY_LIST), message.Payload)
		return nil
	})

	// place

	rpc.Request(messaging.RpcMessage{messaging.LIST, nil})
}

func GetShiftDetail(s socketio.Conn, data interface{}){
	log.Info("Getting shift details!!!")
	utils.TimeTrack(time.Now(), "Get Shift Details")

	log.Info("Finding a shift details")
	payload := data.(map[string]interface{})["payload"]

	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	var shift shifts.Shift
	err = json.Unmarshal(tmp, &shift)

	shift = shifts.GetOneShift(bson.M{"_id": shift.ID})

	for k, v := range shift.Volunteers {
		v.Notifications = nil
		v.CognitoData = nil
		shift.Volunteers[k] = v
	}

	s.Emit(constructSocketID(SHIFT_DETAILS), shift)
}

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
	shift.Manager = user

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

	templateStrings := []string{
		`Hey There is an open shift from {{.StartTime }} to {{.EndTime}} On {{.Date }}. Reply "1" if you would like to pick up this shift. Skip the text messages and download the Flexable app in the Apple App store or the Google play store.`,
		`Hey There is an new shift from {{.StartTime }} to {{.EndTime}} On {{.Date }}. Check it out in your dashboard!`,
	}

	// The New Implementation
	users, err := account.FindAll(bson.M{"profile.company.uuid": shift.Company.UUID})

	for k, v := range templateStrings {
		
	buf, err := CreateTextMessageString(v, shift)

		if err != nil {
			panic(err)
		}

		templateStrings[k] = buf.String()
	}
	for _, user := range users {
		user.Notify(templateStrings, NEW_SHIFT_TITLE, shift.PhoneNumber, shift.Manager)
	}
}

func CloseoutShift(s socketio.Conn, data interface{}){
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

	shift.ClosedOut.Signor = user
	shift.ClosedOut.Closed = true

	shift.Save()


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

func GetAvailableEmployees(s socketio.Conn, data interface{}) {

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

			user := account.User{
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

			user.Upsert(bson.M{"_id": user.ID})
			empList = append(empList, user)

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
