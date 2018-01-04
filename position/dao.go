package position

import (
	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "job"

func Find(query bson.M) *mgo.Query {

	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}

	c := session.DB(database).C(COLLECTION)

	return c.Find(query)
}

func GetAllPositions(query bson.M) (result []Position) {

	err := Find(query).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result

}

func GetOnePosition(query bson.M) (result Position) {

	err := Find(query).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (position *Position) Save() {

	session, database, err := utils.GetMgoSession()

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(database).C(COLLECTION)

	if position.ID == "" {
		position.ID = bson.NewObjectId()
	}

	info, err := c.UpsertId(position.ID, &position)
	log.Info(info)

	if err != nil {
		panic(err)
	}

}
