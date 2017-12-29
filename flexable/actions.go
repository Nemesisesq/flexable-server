package flexable

import (
	"fmt"

	"github.com/nemesisesq/flexable/protobuf"
	"github.com/odknt/go-socket.io"
)

func SetListeners(socket *socketio.Server) {
	for _, i := range messageTypes {

		id := constructSocketID(i.T)
		socket.OnEvent("/", id, i.H)

	}
}
func constructSocketID(payload_type payload.Payload_Type) string {
	return fmt.Sprintf("socket%d", payload_type)

}

//func ( type, payload) = > {
//	socket.emit(`socket${type}`, payload);
//}
