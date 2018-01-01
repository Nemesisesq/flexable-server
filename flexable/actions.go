package flexable

import "github.com/odknt/go-socket.io"

func SetListeners(socket *socketio.Server) {
	for _, i := range messageTypes {

		id := constructSocketID(i.T)
		socket.OnEvent("/", id, i.H)

	}
}

//func ( type, payload) = > {
//	socket.emit(`socket${type}`, payload);
//}
