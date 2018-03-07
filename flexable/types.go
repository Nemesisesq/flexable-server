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

//Employee Constants

const (
	GET_EMPLOYEE_SHIFTS = payload.EmployeePayload_GET_SHIFTS
)

const (
	MANAGER  = "manager"
	EMPLOYEE = "employee"
)

type ProtoBuffer interface {
	String() string
	EnumDescriptor() ([]byte, []int)
}

type SockHandler func(socketio.Conn, interface{}) interface{}
type MessageType struct {
	T ProtoBuffer
	H SockHandler
	N string
}

var messageTypes = []MessageType{
	{OPEN_SHIFTS, OpenShiftHandler, MANAGER},
	{SHIFT_DETAILS, nil, MANAGER},
	{ACCEPT_SHIFT_SUBSTITUE, nil, MANAGER},
	{DENY_SHIFT_SUBSTITUTE, nil, MANAGER},
	{FIND_SHIFT_SUBSTITUTE, FindShiftReplacementHandler, MANAGER},
	{GET_AVAILABLE_EMPLOYEES, nil, MANAGER},
	{GET_JOBS, GetPositions, MANAGER},
	{EMPLOYEE_LIST, GetAvailableEmployees, MANAGER},
	{SELECT_VOLUNTEER, SelectVolunteer, MANAGER},
	{GET_EMPLOYEE_SHIFTS, GetEmployeeShifts, EMPLOYEE},
}

func constructSocketID(payload_type ProtoBuffer) string {
	return fmt.Sprintf("socket%d", payload_type)

}
