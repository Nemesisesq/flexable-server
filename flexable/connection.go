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
		ctx, cancel := context.WithCancel(ctx)
		ctx = context.WithValue(ctx, "cancel", cancel)
		ctx = SetMongoSession(ctx)

		s.SetContext(ctx)

		InitWatchers(s)

		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		ctx := s.Context()
		fmt.Println(ctx)
		//cancel := ctx.Value("cancel").(context.CancelFunc)
		//cancel()
		fmt.Println("meet error:", e)
		fmt.Println("everything cancelled", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		ctx := s.Context().(context.Context)

		cancel := ctx.Value("cancel").(context.CancelFunc)
		cancel()
		fmt.Println("closed and cancelled", msg)
	})
}
func SetMongoSession(i context.Context) context.Context {

	mongodb_uri := os.Getenv("MONGODB_URI")

	dialInfo, err := mgo.ParseURL(mongodb_uri)
	session, err := mgo.Dial(mongodb_uri)

	if err != nil {
		panic(err)
	}

	go func() {

		defer session.Close()

		<-i.Done()
	}()

	ctx := context.WithValue(i, "mgo", session)
	ctx = context.WithValue(ctx, "db", session.DB(dialInfo.Database))
	return ctx
}
