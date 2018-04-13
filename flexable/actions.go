package flexable

import (
	"fmt"

	"github.com/odknt/go-socket.io"
)

func SetListeners(socket *socketio.Server) {
	for _, i := range messageTypes {

		id := constructSocketID(i.T)
		socket.OnEvent(fmt.Sprintf("/%v", i.N), id, i.H)
	}
}
