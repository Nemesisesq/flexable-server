package shifts

import (
	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Find(query bson.M) *mgo.Query {

	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}

	c := session.DB(database).C("shifts")

	return c.Find(query)
}

func GetAllShifts(query bson.M) (result []Shift) {

	err := Find(query).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result

}

func GetOneShift(query bson.M) (result Shift) {

	err := Find(query).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (shift Shift) Save() {

	session, database, err := utils.GetMgoSession()

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(database).C("shifts")

	shift.ID = bson.NewObjectId()

	if shift.ID != "" {
		info, err := c.UpsertId(shift.ID, &shift)
		log.Info(info)

		if err != nil {
			panic(err)
		}
	} else {
		err := c.Insert(&shift)
		if err != nil {
			panic(err)
		}
	}

}
