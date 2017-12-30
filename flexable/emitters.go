package flexable

import "github.com/odknt/go-socket.io"

type SocketContainer struct {
	Socket socketio.Server
}

func InitEmitters(socket interface{}) {

	go func() {
		for {
			select {
			//case msg := Some Channel
			}
		}
	}()

}
