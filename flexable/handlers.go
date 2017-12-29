package flexable

import (
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

	s.Emit("socket0", shift_list, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
	return "hello"
}

func FindShiftReplacementHandler(s socketio.Conn, data interface{}) interface{} {
	shift := data.(shifts.Shift)
	shift.Company = company.Company{"flexable", "123"}
	shift.SmsID = uuid.NewV4().String()

	plivio, err := PlivioClient.NewClient()

	if err != nil {
		panic(err)
	}

	app := plivio.CreateApplication(shift)

	shift.Application = app
	shift.Save()

	number, err := plivio.BuyPhoneNumber(shift)

	shift.PhoneNumber = number

	shift.Save()

	plivio.SendMessages(shift)
	if err != nil {
		panic(err)
	}

	return nil
}
