package employee

import (
	"fmt"

	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "employee"

var ch chan bool

func Find(query bson.M) *mgo.Query {

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

	c := session.DB(database).C(COLLECTION)

	return c.Find(query)
}

func GetAllEmployees(query bson.M) (result []Employee) {

	err := Find(query).All(&result)
	if err != nil {
		panic(err)
	}

	ch <- true

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
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	return result
}

func (employee *Employee) Save() {

	session, database, err := utils.GetMgoSession()

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(database).C(COLLECTION)

	if employee.ID == "" {
		employee.ID = bson.NewObjectId()
	}

	info, err := c.UpsertId(employee.ID, &employee)
	log.Info(info)

	if err != nil {
		panic(err)
	}

}
