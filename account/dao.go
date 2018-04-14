package account

import (
	"encoding/json"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/oxequa/grace"
	"github.com/nemesisesq/flexable/db"
)

func UserRole(r http.Request) (string, interface{}) {
	tmp := map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tmp)
	if err != nil {
		grace.Recover(&err)
	}

	user := GetUser(bson.M{"profile.email": tmp["email"]})
	mapstructure.Decode(tmp, &user.CognitoData)
	if user.ID == "" {
		fmt.Println("setting user id")

		user.ID = bson.NewObjectId()
		user.Profile.FirstName = "First"
		user.Profile.LastName = "Last"
	}
	user.Profile.PhoneNumber = user.CognitoData["phone_number"].(string)
	user.Profile.Email = tmp["email"].(string)
	user.Upsert(bson.M{"_id": user.ID})
	return user.Role, user.Profile
}
func GetUser(query bson.M) User {
	session := db.GetMgoSession()
	db := session.DB(db.FLEXABLE)
	collection := db.C("user")
	user := &User{}
	err := collection.Find(query).One(&user)
	if err != nil {
		grace.Recover(&err)
	}
	return *user
}
func (user *User) Find(query bson.M) *User {
	tmp := GetUser(query)
	*user = tmp
	return user
}
func (user *User) Upsert(query bson.M) {
	session := db.GetMgoSession()
	db := session.DB(db.FLEXABLE)
	collection := db.C("user")

	id, err := collection.Upsert(query, &user)
	if err != nil {
		grace.Recover(&err)
	}
	log.Info(id)
	if err != nil {
		grace.Recover(&err)
	}
}
func SavePushToken(r http.Request) {
	tmp := map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tmp)
	if err != nil {
		defer grace.Recover(&err)

	}
	user := GetUser(bson.M{"profile.email": tmp["email"]})
	user.PushToken = tmp["token"].(string)
	user.Upsert(bson.M{"_id": user.ID})
}

func FindAll(query bson.M) (users []User, err error) {
	session := db.GetMgoSession()
	defer session.Close()
	if err != nil {
		grace.Recover(&err)
	}
	db := session.DB(db.FLEXABLE)
	collection := db.C("user")
	err = collection.Find(query).All(&users)
	if err != nil {
		grace.Recover(&err)
	}
	return users, nil
}
