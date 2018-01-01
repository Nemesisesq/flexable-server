package utils

import (
	"os"

	"gopkg.in/mgo.v2"
)

func GetMgoSession() (string *mgo.Session, database string, err error) {
	mongodb_uri := os.Getenv("MONGODB_URI")

	dialInfo, err := mgo.ParseURL(mongodb_uri)

	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(mongodb_uri)
	//session, err := mgo.DialWithTimeout(mongodb_uri, time.Second*5)
	return session, dialInfo.Database, err
}
