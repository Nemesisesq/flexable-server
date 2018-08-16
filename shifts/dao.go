package shifts

import (
	"time"

	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"github.com/nemesisesq/flexable-server/db"
	"github.com/oxequa/grace"
	"github.com/nemesisesq/flexable-server/utils"
)

//Channel to close conections to Mongo from other methods that get  query
func GetAllShifts(query bson.M) (result []Shift) {
	session := db.GetMgoSession()

	defer session.Close()
	c := session.DB(db.FLEXABLE).C("shifts")
	err := c.Find(query).All(&result)

	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	out := []Shift{}
	layout := "2006-01-02T15:04:05.000Z"
	for _, v := range result {

		end, err := time.Parse(layout, v.RawEndTime)
		if err != nil {
			log.Error(err)
		}

		now := time.Now()
		//log.Info(now.Format(layout))
		//log.Info(now.String())

		//log.Info(now.Hour(), now.Minute())1
		//log.Info(then.Hour(), then.Minute())


		_ = now.Before(end)


		if v.ClosedOut.Closed != true{
			out = append(out, v)
		}

		//	Filter closed sessions

	}

	return out

}

func GetOneShift(query bson.M) (result Shift) {
	utils.TimeTrack(time.Now(), "Get One Shift Query")

	session := db.GetMgoSession()
	defer session.Close()
	c := session.DB(db.FLEXABLE).C("shifts")

	err := c.Find(query).One(&result)
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

func (shift Shift) Save() {
	session := db.GetMgoSession()
	defer session.Close()
	c := session.DB(db.FLEXABLE).C("shifts")
	if shift.ID == "" {
		id := bson.NewObjectId()
		log.Debug(id.String())
		log.Debug(id.Hex())
		log.Debug(id.Machine())
		shift.ID = id
	}
	info, err := c.UpsertId(shift.ID, &shift)
	log.Info(info)
	if err != nil {
		panic(err)
	}
}
