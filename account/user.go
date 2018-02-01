package account

import "gopkg.in/mgo.v2/bson"

type Permission struct {
}

type Group struct {
}

type GeoLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Jobs struct {
}

type User struct {
	ID          bson.ObjectId
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Location    GeoLocation
	Permissions []Permission
	Groups      []Group
	JobHistory  []Jobs
	Auth0ID     string
}
