package flexable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/employee"
	PlivoClient "github.com/nemesisesq/flexable/plivio/client"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func OpenShiftHandler(s socketio.Conn, _ interface{}) interface{} {
	fmt.Print("hello 0")

	shift_list := shifts.GetAllShifts(nil)

	s.Emit(constructSocketID(OPEN_SHIFTS), shift_list, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
	return "hello"
}

func FindShiftReplacementHandler(s socketio.Conn, data interface{}) interface{} {
	payload := data.(map[string]interface{})["payload"]
	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	var shift shifts.Shift
	err = json.Unmarshal(tmp, &shift)

	if err != nil {
		panic(err)
	}
	shift.Company = company.Company{"flexable", "123"}
	shift.SmsID = uuid.NewV4().String()

	plivoClient, err := PlivoClient.NewClient()

	if err != nil {

		panic(err)
	}

	app := plivoClient.CreateApplication(shift)

	shift.Application = app
	shift.Save()
	shift.Name = fmt.Sprint("Rando", time.Now())
	// TODO uncomment for testing buying phone number works
	//number, err := plivio.BuyPhoneNumber(shift)

	response, err := plivoClient.Numbers.Update(
		"16143636301",
		plivo.NumberUpdateParams{
			AppID: app.AppID,
		},
	)

	log.Info(response)
	if err != nil {
		panic(err)
	}

	//shift.PhoneNumber = number
	shift.PhoneNumber = "16143636301"

	shift.Save()

	templateString := `
Hey There is an open shift from {{.StartTime }} to {{.EndTime}}
On {{.Date }}
`

	buf, err := CreateTextMessageString(templateString, shift)
	plivoClient.SendMessages(shift.PhoneNumber, "12165346715<16142881847", buf.String())
	if err != nil {
		panic(err)
	}

	return nil
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
			V            int    `json:"__v"`
			ID           string `json:"_id"`
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
			Job  struct {
				Compensation int    `json:"compensation"`
				Title        string `json:"title"`
			} `json:"job"`
			Location    string `json:"location"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
			SmsID       string `json:"sms_id"`
			Start       string `json:"start"`
			Volunteers  struct {
				Number string `json:"number"`
			} `json:"volunteers"`
		} `json:"shift"`
		Volunteer employee.Employee `json:"volunteer"`
	} `json:"payload"`
}

func SelectVolunteer(s socketio.Conn, data interface{}) interface{} {
	//fmt.Println("Selecting the Volunteer")
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

	err = plivoClient.SendMessages(shift.PhoneNumber, shift.Chosen.Number, buf.String())
	if err != nil {
		panic(err)
	}

	shift.Save()

	log.Info(shift)

	return nil
}
