package account

import (
	"encoding/json"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"github.com/kr/pretty"
	"fmt"
)

var ch = make(chan bool, 1)

func UserRole(r http.Request) (string, interface{}) {
	tmp := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&tmp)

	if err != nil {
		panic(err)
	}

	user := GetUser(bson.M{"email": tmp["email"]})

	mapstructure.Decode(tmp, &user.CognitoData)
	pretty.Println(user)
	if user.ID == "" {
		fmt.Println("setting user id")

		user.ID = bson.NewObjectId()
	}
	user.Email = tmp["email"].(string)

	user.Upsert(bson.M{"email": tmp["email"]})

	return user.Role, user.Profile

}
func GetUser(query bson.M) User {
	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}
	db := session.DB(database)
	collection := db.C("user")
	user := &User{}
	err = collection.Find(query).One(&user)

	//user.CognitoData = tmp

	return *user
}
func (user *User) Find(query bson.M) *User {
	tmp := GetUser(query)
	*user = tmp
	return user
}
func (user *User) Upsert(query bson.M) {
	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}
	db := session.DB(database)
	collection := db.C("user")

	id, err := collection.Upsert(query, &user)
	if err != nil {
		panic(err)
	}
	log.Info(id)
	if err != nil {
		panic(err)
	}
}
func SavePushToken(r http.Request) {
	tmp := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tmp)

	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}

	go func() {
		for {

			select {
			case <-ch:
				session.Close()
			}
		}
	}()

	db := session.DB(database)

	collection := db.C("user")

	collection.Upsert(bson.M{"auth0_id": tmp["auth0_id"]}, bson.M{"token": tmp["auth"]})
	ch <- true
}
