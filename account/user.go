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
	ID          bson.ObjectId          `json:"id" bson:"id"`
	FirstName   string                 `json:"first_name" bson:"first_name"`
	LastName    string                 `json:"last_name" bson:"last_name"`
	Email       string                 `json:"email" bson:"email"`
	PhoneNumber string                 `json:"phone_number" bson:"phone_number"`
	Location    GeoLocation            `json:"location" bson:"location"`
	Permissions []Permission           `json:"permissions" bson:"permissions"`
	Groups      []Group                `json:"groups" bson:"groups"`
	JobHistory  []Jobs                 `json:"job_history" bson:"job_history"`
	Auth0ID     string                 `json:"auth_0_id" bson:"auth_0_id"`
	Role        string                 `json:"role" bson:"role"`
	CognitoData map[string]interface{} `json:"cognito_data" bson:"cognito_data"`
}
