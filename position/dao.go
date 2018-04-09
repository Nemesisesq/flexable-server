package position

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/nemesisesq/flexable/db"
)

const COLLECTION = "job"

var ch chan bool

func Find(query bson.M) *mgo.Query {
	ch = make(chan bool)
	session:= db.GetMgoSession()

	go func() {
		for {
			select {
			case <-ch:
				session.Close()
			}
		}
	}()

	c := session.DB(db.FLEXABLE).C(COLLECTION)

	return c.Find(query)
}

func GetAllPositions(query bson.M) (result []Position) {

	err := Find(query).All(&result)
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	ch <- true
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

	session := db.GetMgoSession()

	defer session.Close()

	c := session.DB(db.FLEXABLE).C(COLLECTION)

	if position.ID == "" {
		position.ID = bson.NewObjectId()
	}

	info, err := c.UpsertId(position.ID, &position)
	log.Info(info)

	if err != nil {
		panic(err)
	}

}
