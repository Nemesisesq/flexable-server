package flexable

import (
	"github.com/nemesisesq/flexable/protobuf"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
)

const (
	OPEN_SHIFTS             = payload.Payload_OPEN_SHIFTS
	SHIFT_DETAILS           = payload.Payload_SHIFT_DETAILS
	ACCEPT_SHIFT_SUBSTITUE  = payload.Payload_ACCEPT_SHIFT_SUBSTITUE
	DENY_SHIFT_SUBSTITUTE   = payload.Payload_DENY_SHIFT_SUBSTITUTE
	FIND_SHIFT_SUBSTITUTE   = payload.Payload_FIND_SHIFT_SUBSTITUTE
	GET_AVAILABLE_EMPLOYEES = payload.Payload_GET_AVAILABLE_EMPLOYEES
	GET_JOBS                = payload.Payload_GET_JOBS
	EMPLOYEE_LIST           = payload.Payload_EMPLOYEE_LIST
)

type SockHandler func(so socketio.Conn, data interface{}) interface{}
type MessageType struct {
	T payload.Payload_Type
	H SockHandler
}

var messageTypes = []MessageType{
	MessageType{OPEN_SHIFTS, shifts.OpenShiftHandler},
	MessageType{SHIFT_DETAILS, nil},
	MessageType{ACCEPT_SHIFT_SUBSTITUE, nil},
	MessageType{DENY_SHIFT_SUBSTITUTE, nil},
	MessageType{FIND_SHIFT_SUBSTITUTE, nil},
	MessageType{GET_AVAILABLE_EMPLOYEES, nil},
	MessageType{GET_JOBS, nil},
	MessageType{EMPLOYEE_LIST, nil},
}