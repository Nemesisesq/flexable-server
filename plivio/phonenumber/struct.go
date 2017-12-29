package phonenumber

import "github.com/plivo/plivo-go/plivo"

type PhoneNumber struct {
	plivo.PhoneNumber
	plivo.PhoneNumberCreateResponse
}
