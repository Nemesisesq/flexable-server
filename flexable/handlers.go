package flexable

import (
	"encoding/json"
	"fmt"

	"github.com/nemesisesq/flexable/company"
	PlivioClient "github.com/nemesisesq/flexable/plivio/client"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func OpenShiftHandler(s socketio.Conn, _ interface{}) interface{} {
	fmt.Print("hello 0")

	shift_list := shifts.GetAllShifts()

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

	plivio, err := PlivioClient.NewClient()

	if err != nil {
		panic(err)
	}

	app := plivio.CreateApplication(shift)

	shift.Application = app
	shift.Save()
	// TODO uncomment for testing buying phone number works
	//number, err := plivio.BuyPhoneNumber(shift)

	//shift.PhoneNumber = number
	shift.PhoneNumber = "16143636301"

	shift.Save()

	plivio.SendMessages(shift)
	if err != nil {
		panic(err)
	}

	return nil
}
