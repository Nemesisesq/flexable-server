package shifts

import "gopkg.in/mgo.v2/bson"

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
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	Name         string        `json:"name"`
	AbsentWorker Worker        `json:"absentWorker" bson:"absentWorker"`
	Job          Job           `json:"job"`
	Location     string        `json:"location"`
	Date         string        `json:"date"`
	StartTime    string        `json:"startTime"`
	EndTime      string        `json:"endTime"`
	Volunteers   []interface{} `json:"volunteers"`
	V            int           `json:"__v"`
}
