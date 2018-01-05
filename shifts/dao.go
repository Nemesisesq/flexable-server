package shifts

import (
	"time"

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

	out := []Shift{}
	for _, v := range result {

		then, err := time.Parse("Mon Jan 2 2006 15:04:05 MST-0700", v.RawEndTime)
		if err != nil {
			log.Error(err)
		}
		now := time.Now()

		//log.Info(now.Hour(), now.Minute())
		//log.Info(then.Hour(), then.Minute())

		if now.Before(then) {
			out = append(out, v)
		}

	}

	return out

}

func GetOneShift(query bson.M) (result Shift) {

	err := Find(query).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func (shift *Shift) Save() {

	session, database, err := utils.GetMgoSession()

	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB(database).C("shifts")

	if shift.ID == "" {
		shift.ID = bson.NewObjectId()
	}

	info, err := c.UpsertId(shift.ID, &shift)
	log.Info(info)

	if err != nil {
		panic(err)
	}

}
