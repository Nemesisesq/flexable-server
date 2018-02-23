package employee

import (
	"github.com/nemesisesq/flexable/position"
	"gopkg.in/mgo.v2/bson"
)

type GeoLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Shiftable interface {
	String()
}

type Employee struct {
	ID       bson.ObjectId     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string            `json:"name" bson:"name"`
	Number   string            `json:"number" bson:"number"`
	Email    string            `json:"email" bson:"email"`
	Location GeoLocation       `json:"location" bson:"location"`
	Position position.Position `json:"position" bson:"position"`
	Schedule []Shiftable       `json:"schedule" bson:"schedule"`
}
