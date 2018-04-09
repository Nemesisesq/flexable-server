package db

import (
	"os"

	"github.com/globalsign/mgo"
)


const FLEXABLE string ="flexable"


var session *mgo.Session

func InitDB() (string, error) {
	mongodb_uri := os.Getenv("MONGODB_URI")

	dialInfo, err := mgo.ParseURL(mongodb_uri)

	if err != nil {
		panic(err)
	}

	session, err = mgo.Dial(mongodb_uri)
	//session, err := mgo.DialWithTimeout(mongodb_uri, time.Second*5)
	return dialInfo.Database, err
}


func GetMgoSession() (*mgo.Session){
	return session.Clone()
}