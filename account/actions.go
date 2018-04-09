package account

import (
	"fmt"
	"github.com/Terminux/exponent-server-sdk-go"
	"github.com/oxequa/grace"
	log "github.com/sirupsen/logrus"
	PlivoClient "github.com/nemesisesq/flexable/plivio/client"

	"github.com/nemesisesq/flexable/shifts"
)

func (u User) Notify(message, title string, shift interface{}) {

	apiRes, apiErr := u.Push(message, title)

	log.Info(apiRes)
	if &apiErr != nil {
		u.Text(message, title, shift.(shifts.Shift))
	}

}

func (u *User) Push(message, title string) (apiRes expo.PushNotificationResult, apiErr expo.PushNotificationError) {

	if expo.IsExpoPushToken(u.PushToken) {
		message := expo.PushMessage{
			To:    u.PushToken,
			Title: title,
			Body:  message,
			Data:  struct{ Value string }{"mydata"}}

		apiRes, apiErr, err := message.Send()
		if err != nil {
			grace.Recover(&err)
		}

		fmt.Println("apiRes:", apiRes)
		fmt.Println("apiErr:", apiErr)

	}
	return apiRes, apiErr

}

func (u *User) Text(message, title string, shift shifts.Shift) (err error) {
	plivoClient, err := PlivoClient.NewClient()

	if err != nil {

		grace.Recover(&err)
	}

	//shift.PhoneNumber = number
	err = plivoClient.SendMessages(shift.PhoneNumber, u.Profile.PhoneNumber, message)
	if err != nil {
		grace.Recover(&err)
	}

	return err
}
