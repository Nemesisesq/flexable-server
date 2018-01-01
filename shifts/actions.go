package shifts

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func UpdateShiftWithSmsID(smsId string, payload map[string]string) bool {

	shift := GetOneShift(bson.M{"sms_id": smsId})

	employee := map[string]interface{}{
		"number": payload["From"],
	}
	if payload["Text"] == "1" {
		shift.Volunteers = append(shift.Volunteers, employee)
		shift.Save()

		return true
		log.Info(payload, shift)
	}

	return false
}
