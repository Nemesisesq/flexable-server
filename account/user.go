package account

import "github.com/globalsign/mgo/bson"

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

type Profile struct {
	Company     string      `json:"company" bson:"company"`
	JobHistory  []Jobs      `json:"job_history" bson:"job_history"`
	PhoneNumber string      `json:"phone_number" bson:"phone_number"`
	Location    GeoLocation `json:"location" bson:"location"`
	FirstName   string      `json:"first_name" bson:"first_name"`
	LastName    string      `json:"last_name" bson:"last_name"`
}

type User struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	Email       string        `json:"email" bson:"email"`
	Permissions []Permission  `json:"permissions" bson:"permissions"`
	Groups      []Group       `json:"groups" bson:"groups"`
	Role        string        `json:"role" bson:"role"`

	Profile Profile `json:"profile" bson:"profile"`

	CognitoData map[string]interface{} `json:"cognito_data" bson:"cognito_data"`
}
