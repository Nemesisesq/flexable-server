package employee

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/nemesisesq/flexable/db"
	"github.com/oxequa/grace"
)

const COLLECTION = "employee"

var ch chan bool

func Find(query bson.M) *mgo.Query {

	session := db.GetMgoSession()

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

func GetAllEmployees(query bson.M) (result []Employee) {
	session := db.GetMgoSession()

	defer session.Close()

	c := session.DB(db.FLEXABLE).C(COLLECTION)
	err := c.Find(query).All(&result)
	if err != nil {
		grace.Recover(&err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	return result

}

func GetOneEmployee(query bson.M) (result Employee) {
	//TODO fix mongo access for employee unpack this find Query

	err := Find(query).One(&result)
	if err != nil {
		grace.Recover(&err)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	return result
}

func (employee *Employee) Save() {

	session := db.GetMgoSession()

	defer session.Close()

	c := session.DB(db.FLEXABLE).C(COLLECTION)

	if employee.ID == "" {
		employee.ID = bson.NewObjectId()
	}

	info, err := c.UpsertId(employee.ID, &employee)
	log.Info(info)

	if err != nil {
		grace.Recover(&err)
	}

}
