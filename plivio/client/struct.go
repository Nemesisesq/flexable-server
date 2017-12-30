package client

import (
	"bytes"
	"fmt"
	application2 "github.com/nemesisesq/flexable/plivio/application"
	"github.com/nemesisesq/flexable/plivio/phonenumber"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/plivo/plivo-go/plivo"
	"github.com/satori/go.uuid"
	"os"
	"regexp"
	"text/template"
)

type Client plivo.Client

func NewClient() (Client, error) {
	client, err := plivo.NewClient("MAMTIXODFLNTK0YWU3NJ", "MDg0YTExMzU4YzdmMWZiZThiZGVjMTM2ODI3NDgx", &plivo.ClientOptions{})
	return Client(*client), err
}

func (c Client) CreateApplication(s shifts.Shift) application2.Application {

	host := os.Getenv("HOST")

	id := uuid.NewV4().String()
	response, err := c.Applications.Create(
		plivo.ApplicationCreateParams{
			AppName:    fmt.Sprint(s.Company.Name, id),
			MessageURL: fmt.Sprintf("%v/sms/incoming/%v", host, s.SmsID),
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

func (c Client) BuyPhoneNumber(s shifts.Shift) (string, error) {
	v, err := c.SearchPhoneNumbers()

	response, err := c.PhoneNumbers.Create(
		v.Number,
		plivo.PhoneNumberCreateParams{
			AppID: s.Application.AppID,
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

func (c Client) SendMessages(s shifts.Shift) error {
	s.Company.GetAvailableWorkers()

	messageTemplate, err := template.New("test").Parse(`
Hey There is an open shift from {{.StartTime }} to {{.EndTime}}
On {{.Date }}
`)

	if err != nil {
		panic(err)
	}
	buf := bytes.Buffer{}
	err = messageTemplate.Execute(&buf, s)
	if err != nil {
		panic(err)
	}
	response, err := c.Messages.Create(
		plivo.MessageCreateParams{
			Src:  s.PhoneNumber,
			Dst:  "12165346715<16142881847",
			Text: buf.String(),
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %#v\n", response)

	return nil
}
