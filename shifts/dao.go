package shifts

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetAllShifts() []Shift {

	session, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("flexable").C("shifts")

	result := []Shift{}
	err = c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result

}

func (shift Shift) Save() {

	session, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("flexable").C("shifts")

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
