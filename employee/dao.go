package employee

import (
	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const COLLECTION = "employee"

func Find(query bson.M) *mgo.Query {

	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}

	c := session.DB(database).C(COLLECTION)

	return c.Find(query)
}

func GetAllEmployees(query bson.M) (result []Employee) {

	err := Find(query).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result

}

func GetOneEmployee(query bson.M) (result Employee) {

	err := Find(query).One(&result)
	if err != nil {
		log.Fatal(err)
	}

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
