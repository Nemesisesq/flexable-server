package account

import (
	"github.com/globalsign/mgo/bson"
	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/position"
)

type Permission struct {
}

type Group struct {
}

type Jobs struct {
}

type Profile struct {
	Company     company.Company `json:"company" bson:"company"`
	JobHistory  []Jobs          `json:"job_history" bson:"job_history"`
	PhoneNumber string          `json:"phone_number" bson:"phone_number"`
	Location    GeoLocation     `json:"location" bson:"location"`
	FirstName   string          `json:"first_name" bson:"first_name"`
	LastName    string          `json:"last_name" bson:"last_name"`
}

type CognitoData struct {
	Sub      string `json:"sub" bson:"sub"`
	EventID  string `json:"event_id" bson:"event_id"`
	TokenUse string `json:"token_use" bson:"token_use"`
	Scope    string `json:"scope" bson:"scope"`
	AuthTime int    `json:"auth_time" bson:"auth_time"`
	Iss      string `json:"iss" bson:"iss"`
	Exp      int    `json:"exp" bson:"exp"`
	Iat      int    `json:"iat" bson:"iat"`
	Jti      string `json:"jti" bson:"jti"`
	ClientID string `json:"client_id" bson:"client_id"`
	Username string `json:"username" bson:"username"`
}

type User struct {
	ID          bson.ObjectId          `json:"_id" bson:"_id"`
	Email       string                 `json:"email" bson:"email"`
	Permissions []Permission           `json:"permissions" bson:"permissions"`
	Groups      []Group                `json:"groups" bson:"groups"`
	Role        string                 `json:"role" bson:"role"`
	Profile     Profile                `json:"profile" bson:"profile"`
	CognitoData map[string]interface{} `json:"cognito_data" bson:"cognito_data"`
	PushToken   string                 `json:"push_token" bson:"push_token"`
	Position position.Position `json:"position" bson:"position"`
}

type GeoLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Shiftable interface {
	String()
}

type Employee struct {
	User
	//ID       bson.ObjectId     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string            `json:"name" bson:"name"`
	//Number   string            `json:"number" bson:"number"`
	//Email    string            `json:"email" bson:"email"`
	Location GeoLocation       `json:"location" bson:"location"`

	Schedule []Shiftable       `json:"schedule" bson:"schedule"`
}
