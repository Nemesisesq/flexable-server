package flexable

import (
	"fmt"

	"github.com/nemesisesq/flexable/protobuf"
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
	SELECT_VOLUNTEER        = payload.Payload_SELECT_VOLUNTEER
)

type SockHandler func(socketio.Conn, interface{}) interface{}
type MessageType struct {
	T payload.Payload_Type
	H SockHandler
}

var messageTypes = []MessageType{
	{OPEN_SHIFTS, OpenShiftHandler},
	{SHIFT_DETAILS, nil},
	{ACCEPT_SHIFT_SUBSTITUE, nil},
	{DENY_SHIFT_SUBSTITUTE, nil},
	{FIND_SHIFT_SUBSTITUTE, FindShiftReplacementHandler},
	{GET_AVAILABLE_EMPLOYEES, nil},
	{GET_JOBS, GetPositions},
	{EMPLOYEE_LIST, GetAvailableEmployees},
	{SELECT_VOLUNTEER, SelectVolunteer},
}

func constructSocketID(payload_type payload.Payload_Type) string {
	return fmt.Sprintf("socket%d", payload_type)

}
