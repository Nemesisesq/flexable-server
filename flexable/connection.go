package flexable

import (
	"context"
	"fmt"
	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
)

func SocketServerConnections(server socketio.Server, namespace string) {
	server.OnConnect(fmt.Sprintf("/%v", namespace), func(s socketio.Conn) error {

		log.Info("Connecting to ", namespace)
		//set context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		ctx = context.WithValue(ctx, "cancel", cancel)
		//ctx = SetMongoSession(ctx)

		s.SetContext(ctx)

		InitWatchers(s)

		log.WithFields(log.Fields{
			"namespace": s.Namespace(),
			"ID":        s.ID(),
		}).Info("connected:")
		return nil
	})

	server.OnError(fmt.Sprintf("/%v", namespace), func(s socketio.Conn, e error) {
		//ctx := s.Context().(context.Context)
		//fmt.Println(ctx)
		//cancel := ctx.Value("cancel").(context.CancelFunc)
		//cancel()
		fmt.Println("meet error:", e)
		fmt.Println("everything cancelled", e)
	})
	server.OnDisconnect(fmt.Sprintf("/%v", namespace), func(s socketio.Conn, msg string) {
		//ctx := s.Context().(context.Context)

		//cancel := ctx.Value("cancel").(context.CancelFunc)
		//cancel()
		//fmt.Println("closed and cancelled", msg)
	})
}

//func SetMongoSession(i context.Context) context.Context {
//
//	mongodb_uri := os.Getenv("MONGODB_URI")
//
//	dialInfo, err := mgo.ParseURL(mongodb_uri)
//	session, err := mgo.Dial(mongodb_uri)
//
//	if err != nil {
//		panic(err)
//	}
//
//	go func() {
//
//		defer session.Close()
//
//		<-i.Done()
//	}()
//
//	ctx := context.WithValue(i, "mgo", session)
//	ctx = context.WithValue(ctx, "db", session.DB(dialInfo.Database))
//	return ctx
//}
