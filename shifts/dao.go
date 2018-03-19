package shifts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func SchemaSetup(c *dgo.Dgraph) {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			sms_id: string @index(exact) .
			name: string .
			absentWorker: uid .
			job: uid .
			location: geo.
			date: uid .
			start: string .
			rawStart: string .
			end: string .
			rawEnd: string .
			volunteers: uid .
			company: uid .
			application: uid
			phone_number: string .
			chosen: uid .
		`,
	})

	if err != nil {
		panic(err)
	}
}

type Root struct {
	Shifts []Shift `json:"shifts"`
}

func GetAllGraphShifts(company_uid string) (result []Shift) {
	client := utils.NewDgraphClient()
	SchemaSetup(client)

	variables := map[string]string{"$company": company_uid}
	query := `query Shifts($company:string){
	    			shifts(func: uid($company)){
						shifts {
						expand(_all_)
						}

					}
		}`
	resp, err := client.NewTxn().QueryWithVars(context.Background(), query, variables)
	if err != nil {
		panic(err)
	}

	var root Root

	fmt.Println(string(resp.Json))
	err = json.Unmarshal(resp.Json, &root)
	return root.Shifts

}

//Channel to close conections to Mongo from other methods that get  query
func GetAllShifts(query bson.M) (result []Shift) {
	session, database, err := utils.GetMgoSession()

	defer session.Close()
	if err != nil {
		panic(err)
	}
	c := session.DB(database).C("shifts")
	err = c.Find(query).All(&result)

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

	session, database, err := utils.GetMgoSession()
	defer session.Close()
	if err != nil {
		panic(err)
	}
	c := session.DB(database).C("shifts")

	err = c.Find(query).One(&result)
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

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
