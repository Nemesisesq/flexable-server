package shifts

import (
	"time"

	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Channel to close conections to Mongo from other methods that get  query
var ch chan bool

func Find(query bson.M) *mgo.Query {
	ch = make(chan bool)
	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ch:
				log.Info("closing session")
				session.Close()
			}
		}
	}()

	c := session.DB(database).C("shifts")

	return c.Find(query)
}

func GetAllShifts(query bson.M) (result []Shift) {

	err := Find(query).All(&result)

	ch <- true

	if err != nil {
		//panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	out := []Shift{}
	for _, v := range result {
		layout := "2006-01-02T15:04:05.000Z"

		then, err := time.Parse(layout, v.RawEndTime)
		if err != nil {
			log.Error(err)
		}

		//log.Info(then.String())
		now := time.Now()
		//log.Info(now.String())

		//log.Info(now.Hour(), now.Minute())
		//log.Info(then.Hour(), then.Minute())

		_ = now.Before(then)
		if true {
			out = append(out, v)
		}

	}

	return out

}

func GetOneShift(query bson.M) (result Shift) {

	err := Find(query).One(&result)
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
