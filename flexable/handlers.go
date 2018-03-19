package flexable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/manveru/faker"
	"github.com/mitchellh/hashstructure"
	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/employee"
	PlivoClient "github.com/nemesisesq/flexable/plivio/client"
	"github.com/nemesisesq/flexable/position"
	payload2 "github.com/nemesisesq/flexable/protobuf"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/nemesisesq/flexable/utils"
	"github.com/odknt/go-socket.io"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func OpenShiftHandler(s socketio.Conn, _ interface{}) interface{} {
	log.Info("Returning openshifts")

	shift_list := shifts.GetAllShifts(nil)
	shift_list := shifts.GetAllGraphShifts(nil)

	s.Emit(constructSocketID(OPEN_SHIFTS), shift_list, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
	return "hello"
}

func FindShiftReplacementHandler(s socketio.Conn, data interface{}) interface{} {
	log.Info("Finding a shift replacement")
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
	shift.Name = fmt.Sprintf("%v : %v per hour", shift.Job.Title, shift.Job.Compensation)
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
On {{.Date }}. Reply "1" if you would like to pick up this shift. Skip the text messages and download the Flexable app in the Apple App store
or the Google play store.
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

	err = plivoClient.SendMessages(shift.PhoneNumber, shift.Chosen.Number, buf.String())
	if err != nil {
		panic(err)
	}

	shift.Save()

	log.Info(shift)

	return nil
}

func GetAvailableEmployees(s socketio.Conn, data interface{}) interface{} {
	log.Info("Getting Available Employees")
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	empList := employee.GetAllEmployees(nil)

	if len(empList) <= 1 {

		for i := 0; i < utils.RandomRange(2, 10); i++ {
			var num string
			bin := utils.RandomRange(1, 2)

			if bin%2 == 0 {
				num = "12165346715"
			} else {
				num = "16142881847"
			}

			x := employee.Employee{
				ID:     bson.NewObjectId(),
				Name:   fake.Name(),
				Number: num,
				Email:  fake.Email(),
				Location: employee.GeoLocation{
					fake.Latitude(),
					fake.Longitude(),
				},
				Position: position.Position{
					ID:           bson.NewObjectId(),
					Title:        fake.JobTitle(),
					Compensation: 10.00,
					Rate:         "hr",
				},
			}

			x.Save()
			empList = append(empList, x)

		}
	}

	fmt.Println("sending employees")

	log.Info(s.Namespace())

	id := constructSocketID(EMPLOYEE_LIST)
	s.Emit(id, empList, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
	return empList
}

func GetPositions(s socketio.Conn, data interface{}) interface{} {
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

	s.Emit(constructSocketID(GET_JOBS), jobs, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})

	return jobs
}

type EmployeeData struct {
}

func GetEmployeeShifts(s socketio.Conn, data interface{}) interface{} {

	payload := data.(map[string]interface{})["payload"]

	if payload == nil {
		return nil
	}
	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	var empl_payload payload2.EmployeePayload
	err = json.Unmarshal(tmp, &empl_payload)

	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Second * 2)
	tickChan := ticker.C
	//companyId := "123"
	var employeeShiftHash uint64

	go func() {
		log.Info("starting employee shift watcher")

		for {
			select {
			case <-tickChan:

				shift_list := shifts.GetAllShifts(bson.M{
					"$and": []bson.M{
						{"volunteers": bson.M{"$size": 0}},
						{"company.uuid": 123},
					}})

				currentEmployee := employee.GetOneEmployee(bson.M{"id": empl_payload.Id})

				//Combine the employee shift list as well as make them unique

				var shift_list2 []employee.Shiftable

				// Loop through and cast shift.Shift to Shiftable
				for _, v := range shift_list {
					var tmp employee.Shiftable = v
					shift_list2 = append(shift_list2, tmp)
				}
				combined := append(shift_list2, currentEmployee.Schedule...)

				unique := []shifts.Shift{}

				for _, v := range combined {
					found := false
					for _, uv := range unique {
						if v.(shifts.Shift).Date == uv.Date {
							found = true
						}
					}

					if !found {
						unique = append(unique, v.(shifts.Shift))
					}

				}

				tmpHash, err := hashstructure.Hash(unique, nil)

				if err != nil {
					panic(err)
				}

				if tmpHash != employeeShiftHash {
					s.Emit(constructSocketID(GET_EMPLOYEE_SHIFTS), unique)
				}

			}

		}

	}()
	return nil
}
