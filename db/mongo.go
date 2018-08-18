package db

import (
	"os"

	"github.com/globalsign/mgo"
)


var FLEXABLE string


var session *mgo.Session

func InitDB() (*mgo.Session, string, error) {
	mongodb_uri := os.Getenv("MONGODB_URI")

	dialInfo, err := mgo.ParseURL(mongodb_uri)


	if err != nil {
		panic(err)
	}


	FLEXABLE = dialInfo.Database

	session, err = mgo.Dial(mongodb_uri)
	return session, dialInfo.Database, err
}


func GetMgoSession() (*mgo.Session){
	return session.Clone()
}