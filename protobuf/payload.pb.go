// Code generated by protoc-gen-go. DO NOT EDIT.
// source: payload.proto

package payload

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Payload_Type int32

const (
	Payload_UNKNOWN                 Payload_Type = 0
	Payload_SHIFT_DETAILS           Payload_Type = 1
	Payload_ACCEPT_SHIFT_SUBSTITUE  Payload_Type = 2
	Payload_DENY_SHIFT_SUBSTITUTE   Payload_Type = 3
	Payload_FIND_SHIFT_SUBSTITUTE   Payload_Type = 4
	Payload_GET_AVAILABLE_EMPLOYEES Payload_Type = 5
	Payload_GET_JOBS                Payload_Type = 6
	Payload_EMPLOYEE_LIST           Payload_Type = 7
	Payload_SELECT_VOLUNTEER        Payload_Type = 8
	Payload_OPEN_SHIFTS             Payload_Type = 9
	Payload_GET_OPEN_SHIFTS         Payload_Type = 10
	Payload_PICK_UP_SHIFT           Payload_Type = 11
	Payload_CALL_OFF_SHIFT          Payload_Type = 12
	Payload_GET_MY_SCHEDULES        Payload_Type = 13
	Payload_USER_PROFILE_UPDATED    Payload_Type = 14
	Payload_SET_PROFILE             Payload_Type = 15
	Payload_UPDATE_NOTIFICATIONS    Payload_Type = 16
	Payload_CLOSEOUT_SHIFT          Payload_Type = 17
	Payload_GET_COMPANY_LIST        Payload_Type = 18
)

var Payload_Type_name = map[int32]string{
	0:  "UNKNOWN",
	1:  "SHIFT_DETAILS",
	2:  "ACCEPT_SHIFT_SUBSTITUE",
	3:  "DENY_SHIFT_SUBSTITUTE",
	4:  "FIND_SHIFT_SUBSTITUTE",
	5:  "GET_AVAILABLE_EMPLOYEES",
	6:  "GET_JOBS",
	7:  "EMPLOYEE_LIST",
	8:  "SELECT_VOLUNTEER",
	9:  "OPEN_SHIFTS",
	10: "GET_OPEN_SHIFTS",
	11: "PICK_UP_SHIFT",
	12: "CALL_OFF_SHIFT",
	13: "GET_MY_SCHEDULES",
	14: "USER_PROFILE_UPDATED",
	15: "SET_PROFILE",
	16: "UPDATE_NOTIFICATIONS",
	17: "CLOSEOUT_SHIFT",
	18: "GET_COMPANY_LIST",
}
var Payload_Type_value = map[string]int32{
	"UNKNOWN":                 0,
	"SHIFT_DETAILS":           1,
	"ACCEPT_SHIFT_SUBSTITUE":  2,
	"DENY_SHIFT_SUBSTITUTE":   3,
	"FIND_SHIFT_SUBSTITUTE":   4,
	"GET_AVAILABLE_EMPLOYEES": 5,
	"GET_JOBS":                6,
	"EMPLOYEE_LIST":           7,
	"SELECT_VOLUNTEER":        8,
	"OPEN_SHIFTS":             9,
	"GET_OPEN_SHIFTS":         10,
	"PICK_UP_SHIFT":           11,
	"CALL_OFF_SHIFT":          12,
	"GET_MY_SCHEDULES":        13,
	"USER_PROFILE_UPDATED":    14,
	"SET_PROFILE":             15,
	"UPDATE_NOTIFICATIONS":    16,
	"CLOSEOUT_SHIFT":          17,
	"GET_COMPANY_LIST":        18,
}

func (x Payload_Type) String() string {
	return proto.EnumName(Payload_Type_name, int32(x))
}
func (Payload_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

type Payload struct {
	Type Payload_Type `protobuf:"varint,1,opt,name=type,enum=Payload_Type" json:"type,omitempty"`
	Data string       `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *Payload) Reset()                    { *m = Payload{} }
func (m *Payload) String() string            { return proto.CompactTextString(m) }
func (*Payload) ProtoMessage()               {}
func (*Payload) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Payload) GetType() Payload_Type {
	if m != nil {
		return m.Type
	}
	return Payload_UNKNOWN
}

func (m *Payload) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*Payload)(nil), "Payload")
	proto.RegisterEnum("Payload_Type", Payload_Type_name, Payload_Type_value)
}

func init() { proto.RegisterFile("payload.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xcf, 0x6e, 0x13, 0x31,
	0x10, 0xc6, 0x49, 0x13, 0x9a, 0x76, 0xd2, 0x4d, 0xa6, 0x43, 0x81, 0x00, 0x97, 0xd2, 0x53, 0x4f,
	0x39, 0xc0, 0x13, 0x38, 0xde, 0x59, 0x6a, 0xea, 0xd8, 0xd6, 0xda, 0x2e, 0xda, 0x93, 0xb5, 0xa8,
	0xb9, 0x21, 0x65, 0x85, 0x72, 0xc9, 0xcb, 0xf0, 0x6e, 0xbc, 0x09, 0xf2, 0x6e, 0x22, 0x21, 0xb8,
	0x59, 0xdf, 0x6f, 0xfc, 0xfd, 0xfc, 0x07, 0x8a, 0xae, 0x3d, 0xfc, 0xd8, 0xb5, 0xcf, 0xab, 0xee,
	0xe7, 0x6e, 0xbf, 0xbb, 0xfb, 0x3d, 0x86, 0xa9, 0x1b, 0x12, 0xfa, 0x08, 0x93, 0xfd, 0xa1, 0xdb,
	0x2e, 0x47, 0xb7, 0xa3, 0xfb, 0xf9, 0xa7, 0x62, 0x75, 0xcc, 0x57, 0xe1, 0xd0, 0x6d, 0xeb, 0x1e,
	0x11, 0xc1, 0xe4, 0xb9, 0xdd, 0xb7, 0xcb, 0xb3, 0xdb, 0xd1, 0xfd, 0x65, 0xdd, 0xaf, 0xef, 0x7e,
	0x8d, 0x61, 0x92, 0x47, 0x68, 0x06, 0xd3, 0x68, 0x1e, 0x8d, 0xfd, 0x66, 0xf0, 0x05, 0x5d, 0x43,
	0xe1, 0x1f, 0x54, 0x15, 0x52, 0xc9, 0x41, 0x28, 0xed, 0x71, 0x44, 0xef, 0xe1, 0x8d, 0x90, 0x92,
	0x5d, 0x48, 0x03, 0xf1, 0x71, 0xed, 0x83, 0x0a, 0x91, 0xf1, 0x8c, 0xde, 0xc1, 0xeb, 0x92, 0x4d,
	0xf3, 0x0f, 0x09, 0x8c, 0xe3, 0x8c, 0x2a, 0x65, 0xca, 0xff, 0xd1, 0x84, 0x3e, 0xc0, 0xdb, 0x2f,
	0x1c, 0x92, 0x78, 0x12, 0x4a, 0x8b, 0xb5, 0xe6, 0xc4, 0x1b, 0xa7, 0x6d, 0xc3, 0xec, 0xf1, 0x25,
	0x5d, 0xc1, 0x45, 0x86, 0x5f, 0xed, 0xda, 0xe3, 0x79, 0x3e, 0xcf, 0x09, 0x26, 0xad, 0x7c, 0xc0,
	0x29, 0xdd, 0x00, 0x7a, 0xd6, 0x2c, 0x43, 0x7a, 0xb2, 0x3a, 0x9a, 0xc0, 0x5c, 0xe3, 0x05, 0x2d,
	0x60, 0x66, 0x1d, 0x9b, 0x41, 0xe7, 0xf1, 0x92, 0x5e, 0xc1, 0x22, 0xf7, 0xfc, 0x1d, 0x42, 0xae,
	0x73, 0x4a, 0x3e, 0xa6, 0xe8, 0x86, 0x0c, 0x67, 0x44, 0x30, 0x97, 0x42, 0xeb, 0x64, 0xab, 0xea,
	0x98, 0x5d, 0x65, 0x45, 0xde, 0xbb, 0x69, 0x92, 0x97, 0x0f, 0x5c, 0x46, 0xcd, 0x1e, 0x0b, 0x5a,
	0xc2, 0x4d, 0xf4, 0x5c, 0x27, 0x57, 0xdb, 0x4a, 0x69, 0x4e, 0xd1, 0x95, 0x22, 0x70, 0x89, 0xf3,
	0x2c, 0xf7, 0x1c, 0x4e, 0x00, 0x17, 0xfd, 0x68, 0x4f, 0x93, 0xb1, 0x41, 0x55, 0x4a, 0x8a, 0xa0,
	0xac, 0xf1, 0x88, 0xbd, 0x4e, 0x5b, 0xcf, 0x36, 0x1e, 0xdf, 0x13, 0xaf, 0x4f, 0x3a, 0x69, 0x37,
	0x4e, 0x98, 0x66, 0xb8, 0x27, 0x7d, 0x3f, 0xef, 0xbf, 0xfa, 0xf3, 0x9f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x95, 0x17, 0xc1, 0x2c, 0xfb, 0x01, 0x00, 0x00,
}
