package twilio

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type SMSPayload struct {
	To   string
	From string
	Body string
}

func SendSMSMessage(payload SMSPayload) (*http.Response, error) {
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
	v.Set("To", payload.To)
	v.Set("From", payload.From)
	v.Set("Body", payload.Body)
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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Info(data["sid"])
		}
	} else {
		log.Error(resp.Status)
	}
	return resp, err
}
