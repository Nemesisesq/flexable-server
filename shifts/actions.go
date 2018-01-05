package shifts

import (
	"github.com/nemesisesq/flexable/employee"
	"github.com/nemesisesq/flexable/plivio/messaging"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func UpdateShiftWithSmsID(smsId string, payload map[string]string) *messaging.Response {

	shift := GetOneShift(bson.M{"sms_id": smsId})

	s := payload["From"]
	ms := bson.M{"number": s}
	chosen := employee.GetOneEmployee(ms)

	if payload["Text"] == "1" {
		log.Info("got a voluteer!!!")
		shift.Volunteers = append(shift.Volunteers, chosen)
		log.Info("saving shift with new volunteer")
		shift.Save()

		response := messaging.Response{
			messaging.Message{
				Dst:   payload["From"],
				Type:  "sms",
				Src:   payload["To"],
				Value: "Thanks for voluteering we will be getting back to you shortly to let you know you got the gig!",
			},
		}

		return &response
	}

	return nil
}
