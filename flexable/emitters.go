package flexable

import "github.com/odknt/go-socket.io"

func InitEmitters(socket socketio.Conn) {

	go func() {
		CheckOpenShifts(socket)
	}()

}
func CheckOpenShifts(socketio.Conn) {

}
