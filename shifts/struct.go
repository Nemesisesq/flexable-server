package shifts

import (
	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/plivio/application"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	Title        string        `json:"title" bson:"title"`
	Compensation float64       `json:"compensation" bson:"compensation"`
	ID           bson.ObjectId `json:"_id" bson:"_id"`
}

type Worker struct {
	Name string        `json:"name" bson:"name"`
	ID   bson.ObjectId `json:"_id" bson:"_id"`
}
type Shift struct {
	ID           bson.ObjectId           `json:"_id" bson:"_id"`
	SmsID        string                  `json:"sms_id" bson:"sms_id"`
	Name         string                  `json:"name"`
	AbsentWorker Worker                  `json:"absentWorker" bson:"absentWorker"`
	Job          Job                     `json:"job"`
	Location     string                  `json:"location"`
	Date         string                  `json:"date"`
	StartTime    string                  `json:"start"`
	EndTime      string                  `json:"end"`
	Volunteers   []interface{}           `json:"volunteers"`
	Company      company.Company         `json:"company"`
	Application  application.Application `json:"application"`
	PhoneNumber  string                  `json:"phone_number"`
	V            int                     `json:"__v"`
}
