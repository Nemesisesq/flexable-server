package shifts

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func UpdateShiftWithSmsID(smsId string, payload map[string]string) bool {

	shift := GetOneShift(bson.M{"sms_id": smsId})

	employee := map[string]interface{}{
		"number": payload["From"],
		"name":   "Employee Bob",
		"ID":     "234",
	}
	if payload["Text"] == "1" {
		log.Info("got a voluteer!!!")
		shift.Volunteers = append(shift.Volunteers, employee)
		log.Info("saving shift with new volunteer")
		shift.Save()

		return true
	}

	return false
}
