package account

import (
	"encoding/json"
	"net/http"

	"github.com/nemesisesq/flexable/utils"
	"gopkg.in/mgo.v2/bson"
)

var ch = make(chan bool, 1)

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
