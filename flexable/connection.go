package flexable

import (
	"context"
	"fmt"
	"os"

	"github.com/odknt/go-socket.io"
	"gopkg.in/mgo.v2"
)

func SocketServerConnections(server socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		//set context
		ctx := context.Background()
		ctx, _ = context.WithCancel(ctx)

		ctx = SetMongoSession(ctx)

		s.SetContext(ctx)
		fmt.Println("Context", s.Context())

		InitWatchers(s)

		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
}
func SetMongoSession(i context.Context) context.Context {

	mongodb_uri := os.Getenv("MONGODB_URI")
	session, err := mgo.Dial(mongodb_uri)

	if err != nil {
		panic(err)
	}

	go func() {

		defer session.Close()

		<-i.Done()
	}()

	ctx := context.WithValue(i, "mgo", session)
	ctx = context.WithValue(ctx, "db", session.DB("flexable"))
	return ctx
}
