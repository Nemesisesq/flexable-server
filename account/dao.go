package account

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/nemesisesq/flexable/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var ch = make(chan bool, 1)

func SchemaSetup(c *dgo.Dgraph) {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			name: string @index(term) .
			uuid: string @index(exact) .
			first_name: string .
			last_name: string .
			email: string @index(exact) .
			phone_number: string .
			location: uid .
			permission: string .
			group: uid .
			job: uid .
			role: string .
			cognito_data: uid .

		`,
	})

	if err != nil {
		panic(err)
	}
}
func UserRoleGraph(r http.Request) string {
	client := utils.NewDgraphClient()
	SchemaSetup(client)
	user := User{}
	//utils.PrintBody(&r)
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	root := RetrieveUserByEmail(user, err, client, r)

	if len(root.Users) == 0 {
		user.Role = "employee"

		mu := &api.Mutation{
			CommitNow: true,
		}
		pb, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}

		mu.SetJson = pb
		assigned, err := client.NewTxn().Mutate(r.Context(), mu)
		if err != nil {
			log.Fatal(err)
		}

		log.Info(assigned)

		return "employee"

	} else {

	}

	return root.Users[0].Role
}

type Root struct {
	Users []User `json:"user"`
}

func RetrieveUserByEmail(user User, err error, client *dgo.Dgraph, r http.Request) Root {
	//	query client
	variables := map[string]string{"$email": user.Email}
	query := `query User($email:string){
	    			user(func: eq(email, $email)){
	    				name
	    				email
	    				role

					}
		}`
	resp, err := client.NewTxn().QueryWithVars(r.Context(), query, variables)
	if err != nil {
		panic(err)
	}

	var root Root
	err = json.Unmarshal(resp.Json, &root)
	return root
}

func UserRole(r http.Request) string {
	tmp := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&tmp)

	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	db := session.DB(database)

	collection := db.C("user")

	user := &User{}

	if err := collection.Find(bson.M{"email": tmp["email"]}).One(&user); err != nil {
		id, err := collection.Upsert(bson.M{"email": tmp["email"]}, tmp)

		if err != nil {
			panic(err)
		}

		log.Info(id)
	}

	if err != nil {
		panic(err)
	}

	return user.Role

}
func SavePushToken(r http.Request) {
	tmp := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tmp)

	session, database, err := utils.GetMgoSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB(database)

	collection := db.C("user")

	collection.Upsert(bson.M{"auth0_id": tmp["auth0_id"]}, bson.M{"token": tmp["auth"]})
	ch <- true
}
