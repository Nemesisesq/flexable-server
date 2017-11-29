package socket

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/manveru/faker"
	payload2 "github.com/nemesisesq/flexable/protobuf"
)

const (
	OPEN_SHIFTS = iota
)

func HandleMessage(out chan []byte, in []byte) {

	payload := &payload2.Payload{}

	proto.Unmarshal(in, payload)

	switch payload.Type {
	case payload2.Payload_OPEN_SHIFTS:
		out <- openShifts(payload)
	}

}

type Location interface{}

type Company struct {
	Name     string
	Location Location
}

type Name struct {
	FirstName string
	LastName  string
	NickNames []string
}

type PhoneNumber struct {
	AreaCode   int
	Prefix     int
	LineNumber int
	Extension  int
	Temp       string
}

type User struct {
	Name        Name
	Email       string
	PhoneNumber PhoneNumber
}

type Manager struct {
	User User
}

type Employee struct {
	User User
}

type ShiftSite struct {
	Name     string
	Location Location
	Company  Company
	Manager  []Manager
}

type Shift struct {
	Employee  Employee
	Manager   Manager
	ShiftSite ShiftSite
	DateTime  time.Time
}

func openShifts(p *payload2.Payload) (res []byte) {
	fake, _ := faker.New("en")

	mgr := Manager{
		User{
			Name{
				fake.FirstName(),
				fake.LastName(),
				[]string{fake.JobTitle()},
			},
			fake.LastName(),
			PhoneNumber{
				Temp: fake.PhoneNumber(),
			},
		},
	}

	comp := Company{
		fake.CompanyName(),
		fake.StreetAddress(),
	}

	shifts := []Shift{}

	for i := 0; i < 4; i++ {
		s := Shift{
			Employee: Employee{
				User: User{
					Name{
						fake.FirstName(),
						fake.LastName(),
						[]string{fake.JobTitle()},
					},
					fake.LastName(),
					PhoneNumber{
						Temp: fake.PhoneNumber(),
					},
				},
			},
			Manager: mgr,
			ShiftSite: ShiftSite{
				fake.JobTitle(),
				fake.StreetAddress(),
				comp,
				[]Manager{mgr},
			},
			DateTime: time.Now().Add(time.Hour * 2),
		}
		shifts = append(shifts, s)
	}

	f, _ := json.Marshal(shifts)
	open_shifts := payload2.Payload{
		payload2.Payload_OPEN_SHIFTS,
		string(f),
	}

	res, err := proto.Marshal(&open_shifts)

	if err != nil {
		panic(err)
	}
	return res
}
