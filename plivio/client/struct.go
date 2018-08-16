package client

import (
	"fmt"
	"os"
	"regexp"

	application2 "github.com/nemesisesq/flexable-server/plivio/application"
	"github.com/nemesisesq/flexable-server/plivio/phonenumber"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	"errors"
	"github.com/oxequa/grace"
)

type Client plivo.Client

func NewClient() (Client, error) {
	client, err := plivo.NewClient("MAMTIXODFLNTK0YWU3NJ", "MDg0YTExMzU4YzdmMWZiZThiZGVjMTM2ODI3NDgx", &plivo.ClientOptions{})
	return Client(*client), err
}

func (c Client) CreateApplication(companyName string, smsId interface{} ) application2.Application {

	host := os.Getenv("HOST")

	if &host == nil {
		panic(errors.New("You need to set a host in the environment variables"))
	}

	id := uuid.NewV4().String()
	response, err := c.Applications.Create(
		plivo.ApplicationCreateParams{
			AppName:    fmt.Sprint(companyName, id),
			MessageURL: fmt.Sprintf("%v/sms/incoming/%v", host, smsId),
		},
	)
	if err != nil {
		panic(err)
	}
	//TODO create new application in the database
	application := &application2.Application{
		response.Message,
		response.AppID,
		response.ApiID,
	}

	application.Save()

	return *application
}

func (c Client) BuyPhoneNumber(applicationId string) (string, error) {
	v, err := c.SearchPhoneNumbers()

	response, err := c.PhoneNumbers.Create(
		v.Number,
		plivo.PhoneNumberCreateParams{
			AppID: applicationId,
		},
	)
	if err != nil {
		panic(err)
	}

	if response.Status == "fulfilled" {
		pn := phonenumber.PhoneNumber{*v, *response}
		pn.Save()

		return pn.Number, nil
	}
	fmt.Printf("Response: %#v\n", response)

	return "", err
}

func (c Client) SearchPhoneNumbers() (*plivo.PhoneNumber, error) {
	phoneNumberList, err := c.PhoneNumbers.List(
		plivo.PhoneNumberListParams{
			CountryISO: "US",
			Pattern:    "614",
			Services:   "sms",
		},
	)
	if err != nil {
		panic(err)
	}
	for _, v := range phoneNumberList.Objects {
		matched, err := regexp.MatchString("1614", v.Number)
		if err != nil {
			panic(err)
		}
		if matched {
			return v, nil
		}
	}
	return nil, err
}

func (c Client) SendMessages(from, to, message string) error {

	response, err := c.Messages.Create(
		plivo.MessageCreateParams{
			Src:  from,
			Dst:  to,
			Text: message,
		},
	)
	if err != nil {
		grace.Recover(&err)
	}
	fmt.Printf("Response: %#v\n", response)

	return nil
}
