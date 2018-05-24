package account

import (
	"fmt"
	"github.com/Terminux/exponent-server-sdk-go"
	"github.com/oxequa/grace"
	log "github.com/sirupsen/logrus"
	PlivoClient "github.com/nemesisesq/flexable/plivio/client"
	"github.com/globalsign/mgo/bson"
	"time"
)

func (u User) Notify(messages []string, title string, to string, sender User) {

	if u.Profile.Available {

		apiRes, apiErr := u.Push(messages[1], title)

		log.Info(apiRes, "apiRes")
		log.Info(apiErr, "apiErr")
		if apiErr.Code != "" {
			//TODO
			u.Text(messages[0], title, to)
		}
	}

	n := Notification{bson.NewObjectId(),sender, time.Now(), messages[1], false}

	u.Notifications = append(u.Notifications, n)

	u.Save()

}

func (u *User) Push(message, title string) (apiRes expo.PushNotificationResult, apiErr expo.PushNotificationError) {
	if expo.IsExpoPushToken(u.PushToken) {
		log.Debug("Sending push")
		message := expo.PushMessage{
			To:       u.PushToken,
			Title:    title,
			Body:     message,
			Data:     struct{ Value string `json:"value"` }{message},
			TTL:      300,
			Priority: "high",
		}

		apiRes, apiErr, err := message.Send()
		if err != nil {
			grace.Recover(&err)
		}

		fmt.Println("apiRes:", apiRes)
		fmt.Println("apiErr:", apiErr)

	}
	return apiRes, apiErr

}

func (u *User) Text(message, title string, from string) (err error) {
	plivoClient, err := PlivoClient.NewClient()

	if err != nil {

		grace.Recover(&err)
	}

	//shift.PhoneNumber = number
	err = plivoClient.SendMessages(from, u.Profile.PhoneNumber, message)
	if err != nil {
		grace.Recover(&err)
	}

	return err
}
