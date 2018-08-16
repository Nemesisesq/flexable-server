package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
	"github.com/nemesisesq/flexable-server/db"
	"github.com/oxequa/grace"
	log "github.com/sirupsen/logrus"
)

type Notification struct {
	ID      bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Sender bson.ObjectId `json:"sender" bson:"sender"`
	Date    time.Time     `json:"date" bson:"date"`
	Message string        `json:"message" bson:"message"`
	Read    bool          `json:"read" bson:"read"`
	UserId  bson.ObjectId `json:"user_id" bson:"user_id"`
}
type Notifications []Notification

func (n *Notification) Upsert(query bson.M) {
	session := db.GetMgoSession()
	db := session.DB(db.FLEXABLE)
	collection := db.C("notifications")

	id, err := collection.Upsert(query, &n)
	if err != nil {
		grace.Recover(&err)
	}
	log.Info(id)
	if err != nil {
		grace.Recover(&err)
	}
}

func (n *Notification) Save() {
	if n.ID == "" {
		n.ID = bson.NewObjectId()
	}
	n.Upsert(bson.M{"_id": n.ID})
}


func (n Notifications) FindAll(query bson.M) (err error) {
	session := db.GetMgoSession()
	defer session.Close()
	if err != nil {
		grace.Recover(&err)
	}
	db := session.DB(db.FLEXABLE)
	collection := db.C("notifications")
	err = collection.Find(query).All(&n)
	if err != nil {
		grace.Recover(&err)
	}
	return err
}
