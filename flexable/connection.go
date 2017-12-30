package flexable

import (
	"fmt"

	"github.com/odknt/go-socket.io"
)

func SocketServerConnections(server socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")

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
