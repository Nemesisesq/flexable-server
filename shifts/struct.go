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
	AbsentWorker account.User            `json:"absentWorker" bson:"absentWorker"`
	Job          position.Position       `json:"job"`
	Location     string                  `json:"location"`
	Date         string                  `json:"date"`
	StartTime    string                  `json:"start"`
	RawStartTime string                  `json:"rawStart"`
	EndTime      string                  `json:"end"`
	RawEndTime   string                  `json:"rawEnd"`
	Volunteers   []account.User          `json:"volunteers"`
	Company      company.Company         `json:"company"`
	Application  application.Application `json:"application"`
	PhoneNumber  string                  `json:"phone_number"`
	Chosen       account.User            `json:"chosen"`
	Manager      account.User            `json:"manager"`
	V            int                     `json:"__v"`
	ClosedOut struct {
		Signor account.User `json:"signor"`
		Closed bool         `json:"closed"`
	} `json:"closed_out"`
	Details struct {
		PositionName    string   `json:"position_name"`
		SkillsRequested []string `json:"skills_requested"`
		Address         Address  `json:"address"`
		Contact struct {
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		} `json:"contact"`
		Pay string `json:"pay"`
	} `json:"details"`
}

type SkinnyShift struct {
	ID         bson.ObjectId     `json:"_id,omitempty" bson:"_id,omitempty"`
	Job        position.Position `json:"job"`
	Date       string            `json:"date"`
	StartTime  string            `json:"start"`
	EndTime    string            `json:"end"`
	Volunteers int               `json:"volunteers"`
	Chosen     bool              `json:"chosen"`
}

type Address struct {
	Name    string `json:"name"`
	Addr1   string `json:"addr_1"`
	Addr2   string `json:"addr_2"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
}

func (s Shift) String() {
	fmt.Println("I'm an Employee")
}
