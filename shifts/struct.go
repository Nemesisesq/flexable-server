package shifts

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/plivio/application"
	"github.com/nemesisesq/flexable/position"
	"github.com/nemesisesq/flexable/account"
)

type Job struct {
	Title        string        `json:"title" bson:"title"`
	Compensation float64       `json:"compensation" bson:"compensation"`
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
}

type Worker struct {
	Name string        `json:"name" bson:"name"`
	ID   bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
}
type Shift struct {
	ID           bson.ObjectId           `json:"_id,omitempty" bson:"_id,omitempty"`
	SmsID        string                  `json:"sms_id" bson:"sms_id"`
	Name         string                  `json:"name"`
	AbsentWorker account.User       `json:"absentWorker" bson:"absentWorker"`
	Job          position.Position       `json:"job"`
	Location     string                  `json:"location"`
	Date         string                  `json:"date"`
	StartTime    string                  `json:"start"`
	RawStartTime string                  `json:"rawStart"`
	EndTime      string                  `json:"end"`
	RawEndTime   string                  `json:"rawEnd"`
	Volunteers   []account.User     `json:"volunteers"`
	Company      company.Company         `json:"company"`
	Application  application.Application `json:"application"`
	PhoneNumber  string                  `json:"phone_number"`
	Chosen       account.User       `json:"chosen"`
	V            int                     `json:"__v"`
}

func (s Shift) String() {
	fmt.Println("I'm an Employee")
}
