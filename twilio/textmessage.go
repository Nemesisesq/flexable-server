package twilio

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func FindShiftReplacement(data map[string]interface{}) (*http.Response, error) {
	accountSid := "AC8babac161b27ec214bed203884635819"
	authToken := "5c575b32cf3208e7a86e849fd0cd697b"
	//callSid := "PNbf2d127871ca9856d3d06e700edbf3a1"

	const urlTemplate = "https://api.twilio.com/2010-04-01/Accounts/{{.AccountSid}}/Messages"

	data := map[string]interface{}{
		"AccountSid": accountSid,
	}

	t := template.Must(template.New("url").Parse(urlTemplate))

	buf := &bytes.Buffer{}

	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}

	urlStr := buf.String()

	v := url.Values{}
	v.Set("To")
	////logrus.Info(toNum)
	v.Set("From", "+12164506822")
	call_in_number := fmt.Sprintf("%v/twiml", os.Getenv("SELF_URL"))
	////logrus.Info(call_in_number)
	v.Set("Url", call_in_number)
	rb := *strings.NewReader(v.Encode())
	// Create Client
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// make request
	resp, err := client.Do(req)
	return resp, err
}
