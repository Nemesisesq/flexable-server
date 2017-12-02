package socket

import (
	"time"

	"github.com/manveru/faker"
	payload2 "github.com/nemesisesq/flexable/protobuf"
	log "github.com/sirupsen/logrus"
)

//func HandleMessage(in []byte) []byte {
//
//	payload := &payload2.Payload{}
//
//	err := proto.Unmarshal(in, payload)
//
//	if err != nil {
//		panic(err)
//	}
//	defer func() {
//		if r := recover(); r != nil {
//			log.Println("Recovered in f", r)
//		}
//	}()
//
//	switch payload.Type {
//	case payload2.Payload_OPEN_SHIFTS:
//		return OpenShifts(payload)
//
//	default:
//		return nil
//	}
//
//	return nil
//
//}

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

func OpenShifts(p *payload2.Payload) (res interface{}) {
	fake, err := faker.New("en")

	if err != nil {
		panic(err)
	}

	mgr := Manager{
		User{
			Name{
				fake.FirstName(),
				fake.LastName(),
				[]string{fake.JobTitle()},
			},
			fake.Email(),
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
					fake.Email(),
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

	//f, err := json.Marshal(shifts)
	if err != nil {
		panic(err)
	}
	//open_shifts := payload2.Payload{
	//	payload2.Payload_OPEN_SHIFTS,
	//	f,
	//}
	//
	//res, err = proto.Marshal(&open_shifts)

	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	return shifts
}
