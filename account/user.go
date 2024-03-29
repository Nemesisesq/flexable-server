package account

import (
	"github.com/globalsign/mgo/bson"
	"github.com/nemesisesq/flexable-server/company"
	"github.com/nemesisesq/flexable-server/position"
	"time"
)

type Permission struct {
}

type Group struct {
}

type Job struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	Start    string `json:"startDate"`
	End      string `json:"endDate"`
	Location string `json:"location"`
}

type Profile struct {
	Email       string          `json:"email" bson:"email"`
	Company     company.Company `json:"company" bson:"company"`
	JobHistory  []Job           `json:"job_history" bson:"job_history"`
	PhoneNumber string          `json:"phone_number" bson:"phone_number"`
	Location    GeoLocation     `json:"location" bson:"location"`
	FirstName   string          `json:"first_name" bson:"first_name"`
	LastName    string          `json:"last_name" bson:"last_name"`
	ImageUrl    string          `json:"image_url" bson:"image_url"`
	AvailableAt int64           `json:"available_at" bson:"available_at"`
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
	ID            bson.ObjectId          `json:"_id,omitempty" bson:"_id,omitempty"`
	Email         string                 `json:"email" bson:"email"`
	Permissions   []Permission           `json:"permissions" bson:"permissions"`
	Groups        []Group                `json:"groups" bson:"groups"`
	Role          string                 `json:"role" bson:"role"`
	Profile       Profile                `json:"profile" bson:"profile"`
	CognitoData   map[string]interface{} `json:"cognito_data" bson:"cognito_data"`
	PushToken     string                 `json:"push_token" bson:"push_token"`
	Position      position.Position      `json:"position" bson:"position"`
	Notifications []Notification         `json:"notifications" bson:"notifications"`
	IsAdmin       bool                   `json:"is_admin" bson:"is_admin"`
	Rating        int                    `json:"rating" bson:"rating"`
}

type Notification struct {
	ID      bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Sender  User          `json:"sender" bson:"sender"`
	Date    time.Time     `json:"date" bson:"date"`
	Message string        `json:"message" bson:"message"`
	Read    bool          `json:"read" bson:"read"`
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
	Name string `json:"name" bson:"name"`
	//Number   string            `json:"number" bson:"number"`
	//Email    string            `json:"email" bson:"email"`
	Location GeoLocation `json:"location" bson:"location"`

	Schedule []Shiftable `json:"schedule" bson:"schedule"`
}
