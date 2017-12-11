package user

import (
	"bytes"
	"text/template"
)

type User struct {
	Name        Name        `json:"name"`
	Email       string      `json:"email"`
	PhoneNumber PhoneNumber `json:"phone_number"`
}
type PhoneNumber struct {
	AreaCode   int    `json:"area_code"`
	Prefix     int    `json:"prefix"`
	LineNumber int    `json:"line_number"`
	Extension  int    `json:"extension"`
	Temp       string `json:"temp"`
}

func (p PhoneNumber) Number() string {
	const numberTemplate = "{{.AreaCode}}{{.Prefix}}{{.LineNumber}}"

	pn := template.Must(template.New("phone_nubmer").Parse(numberTemplate))
	buf := &bytes.Buffer{}

	if err := pn.Execute(buf, pn); err != nil {
		panic(err)
	}

	return buf.String()
}

type Name struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	NickNames []string `json:"nick_names"`
}
