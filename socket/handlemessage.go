package socket

import (
	"github.com/manveru/faker"
	payload2 "github.com/nemesisesq/flexable/protobuf"
)

type Location interface{}

var fake, _ = faker.New("en")

func OpenShifts(p *payload2.Payload) (res interface{}) {
	return nil
}
