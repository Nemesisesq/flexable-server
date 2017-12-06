package database

import "gopkg.in/mgo.v2"

func GetNumbers(busId string) []string {
	url := "mongo://localhost:27017"
	session, err := mgo.Dial(url)

	result := map[string]interface{}
	c := session.DB("flexable").C("replacement_candidates")
	err := c.Find(bson.M{}).All(&result)
}
