package shifts

import (
	"fmt"

	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
)

func OpenShiftHandler(s socketio.Conn, data interface{}) interface{} {
	fmt.Print("hello 0")

	shift_list := GetAllShifts()

	s.Emit("socket0", shift_list, func(so socketio.Conn, data string) {
		log.Println("Client ACK with data: ", data)
	})
	return "hello"
}
