package flexable

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nemesisesq/flexable/company"
	PlivoClient "github.com/nemesisesq/flexable/plivio/client"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
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

	//log.Info(response)
	if err != nil {
		panic(err)
	}

	//shift.PhoneNumber = number
	shift.PhoneNumber = "16143636301"

	shift.Save()

	plivoClient.SendMessages(shift)
	if err != nil {
		panic(err)
	}

	return nil
}
