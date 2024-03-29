package flexable

import (
	"fmt"

	"github.com/odknt/go-socket.io"
)

func SetListeners(server *socketio.Server) {
	for _, i := range messageTypes {

		id := constructSocketID(i.T)
		server.OnEvent(fmt.Sprintf("/%v", i.N), id, i.H)
	}
}
