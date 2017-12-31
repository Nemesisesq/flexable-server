package flexable

import (
	"context"
	"fmt"

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

		InitEmitters(s)

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

	session, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	go func() {

		defer session.Close()

		<-i.Done()
	}()

	ctx := context.WithValue(i, "mgo", session)
	ctx := context.WithValue(i, "db", session.DB("flexable"))
	return ctx
}
