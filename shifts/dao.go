package shifts

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
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
