package utils

import (
	"context"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewDgraphClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	conn, err := grpc.Dial("ec2-18-218-72-133.us-east-2.compute.amazonaws.com:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(conn),
	)
}

func setup(c *dgo.Dgraph) {
	// Install a schema into dgraph. Accounts have a `name` and a `balance`.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			name: string @index(term) .
			balance: int .
		`,
	})

	if err != nil {
		panic(err)
	}
}
