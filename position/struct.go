package position

import "gopkg.in/mgo.v2/bson"

type Position struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title        string        `json:"title" bson:"title"`
	Compensation float32       `json:"compensation" bson:"compensation"`
	Rate         string        `json:"rate" bson:"rate"`
}
