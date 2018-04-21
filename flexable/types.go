package flexable

import (
	"fmt"

	"github.com/nemesisesq/flexable/protobuf"
	"github.com/odknt/go-socket.io"
	"net/http"
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
	UPDATE_NOTIFICATIONS    = payload.Payload_UPDATE_NOTIFICATIONS
)

//Employee Constants

const (
	GET_OPEN_SHIFTS      = payload.Payload_GET_OPEN_SHIFTS
	PICK_UP_SHIFT        = payload.Payload_PICK_UP_SHIFT
	CALL_OFF_SHIFT       = payload.Payload_CALL_OFF_SHIFT
	USER_PROFILE_UPDATED = payload.Payload_USER_PROFILE_UPDATED
	SET_PROFILE          = payload.Payload_SET_PROFILE
	//	todo fix

)

const (
	MANAGER  = "flexable"
	EMPLOYEE = "flexable"
)

type ProtoBuffer interface {
	String() string
	EnumDescriptor() ([]byte, []int)
}

type SockHandler func(socketio.Conn, interface{}) interface{}
type SocketMessageType struct {
	T ProtoBuffer
	H SockHandler
	N string
}

type HTTPMessageType struct {
	T ProtoBuffer
	H func(http.ResponseWriter, *http.Request)
	N string
}

var messageTypes = []SocketMessageType{
	{OPEN_SHIFTS, OpenShiftHandler, MANAGER},
	{SHIFT_DETAILS, nil, MANAGER},
	{ACCEPT_SHIFT_SUBSTITUE, nil, MANAGER},
	{DENY_SHIFT_SUBSTITUTE, nil, MANAGER},
	{FIND_SHIFT_SUBSTITUTE, FindShiftReplacementHandler, MANAGER},
	{GET_AVAILABLE_EMPLOYEES, nil, MANAGER},
	{GET_JOBS, GetPositions, MANAGER},
	{EMPLOYEE_LIST, GetAvailableEmployees, MANAGER},
	//{SELECT_VOLUNTEER, SelectVolunteer, MANAGER},
	//{GET_MY_SHIFTS, GetEmployeeShifts, EMPLOYEE},
	{GET_OPEN_SHIFTS, GetOpenShifts, EMPLOYEE},
	{PICK_UP_SHIFT, PickUpShift, EMPLOYEE},
	{CALL_OFF_SHIFT, CallOfShift, EMPLOYEE},
	{SET_PROFILE, UpdateProfile, EMPLOYEE},
	{UPDATE_NOTIFICATIONS, UpdateNotifications, EMPLOYEE},
}

func constructSocketID(payload_type ProtoBuffer) string {
	return fmt.Sprintf("socket%d", payload_type)

}


var HttpTypes = []HTTPMessageType{
	{GET_JOBS, GetPositionsHttp, MANAGER},
}