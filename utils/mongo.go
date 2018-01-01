package utils

import (
	"os"

	"gopkg.in/mgo.v2"
)

func GetMgoSession() (string *mgo.Session, database string, err error) {
	mongodb_uri := os.Getenv("MONGODB_URI")

	dialInfo, err := mgo.ParseURL(mongodb_uri)

	session, err := mgo.Dial(mongodb_uri)
	return session, dialInfo.Database, err
}
