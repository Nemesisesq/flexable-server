package shifts

import (
	"github.com/nemesisesq/flexable/employee"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func UpdateShiftWithSmsID(smsId string, payload map[string]string) bool {

	shift := GetOneShift(bson.M{"sms_id": smsId})

	s := payload["From"]
	ms := bson.M{"number": s}
	chosen := employee.GetOneEmployee(ms)

	if payload["Text"] == "1" {
		log.Info("got a voluteer!!!")
		shift.Volunteers = append(shift.Volunteers, chosen)
		log.Info("saving shift with new volunteer")
		shift.Save()

		return true
	}

	return false
}
