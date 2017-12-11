package employees

import (
	"bytes"
	"net/http"
	"text/template"

	"github.com/nemesisesq/flexable/twilio"
	"github.com/nemesisesq/flexable/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Employee struct {
	User      user.User `json:"user"`
	Available bool      `json:"available"`
}

func (e Employee) Contact(shiftDetails map[string]interface{}) (*http.Response, error) {

	const bodyTemplate = ` Hey {{.Name }} i need to fill a shift on {{.ShiftDate}}
from {{.ShiftStartTime }} to {{.ShiftEndTime }} would you like to fill this shift?
reply Yes or No to this message
`
	t := template.Must(template.New("url").Parse(bodyTemplate))

	buf := &bytes.Buffer{}

	if err := t.Execute(buf, shiftDetails); err != nil {
		panic(err)
	}

	messagePayload := twilio.SMSPayload{
		To:   e.User.PhoneNumber.Number(),
		From: "+12164506822",
		Body: buf.String(),
	}

	resp, err := twilio.SendSMSMessage(messagePayload)

	return resp, err
}

type Employees []Employee

func (es Employees) GetAvailable() error {
	url := "mongo://localhost:27017"
	session, err := mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	result := Employees{}
	c := session.DB("flexable").C("replacement_candidates")
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}

	//Find Eligible replacements

	eli := Employees{}
	for _, v := range result {
		if v.Available {
			eli = append(eli, v)
		}
	}

	//TODO Filter employees that are working a shift during vacant shift

	//

	return nil
}
