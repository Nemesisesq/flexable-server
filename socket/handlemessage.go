package socket

import (
	"math/rand"
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
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

type Name struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	NickNames []string `json:"nick_names"`
}

type PhoneNumber struct {
	AreaCode   int    `json:"area_code"`
	Prefix     int    `json:"prefix"`
	LineNumber int    `json:"line_number"`
	Extension  int    `json:"extension"`
	Temp       string `json:"temp"`
}

type User struct {
	Name        Name        `json:"name"`
	Email       string      `json:"email"`
	PhoneNumber PhoneNumber `json:"phone_number"`
}

type Manager struct {
	User User `json:"user"`
}

type Employee struct {
	User User `json:"user"`
}

type ShiftSite struct {
	Name     string    `json:"name"`
	Location Location  `json:"location"`
	Company  Company   `json:"company"`
	Manager  []Manager `json:"manager"`
}

type Shift struct {
	Employee     Employee   `json:"employee"`
	Manager      Manager    `json:"manager"`
	ShiftSite    ShiftSite  `json:"shift_site"`
	DateTime     time.Time  `json:"date_time"`
	Replacements []Employee `json:"replacements"`
	Chosen       Employee   `json:"chosen"`
}

var fake, _ = faker.New("en")

func OpenShifts(p *payload2.Payload) (res interface{}) {

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

		employee := generateEmmployee()
		s := Shift{
			Employee: employee,
			Manager:  mgr,
			ShiftSite: ShiftSite{
				fake.JobTitle(),
				fake.StreetAddress(),
				comp,
				[]Manager{mgr},
			},
			DateTime:     time.Now().Add(time.Hour * 2),
			Replacements: getReplacements(mgr),
		}
		shifts = append(shifts, s)

	}

	//f, err := json.Marshal(shifts)
	//open_shifts := payload2.Payload{
	//	payload2.Payload_OPEN_SHIFTS,
	//	f,
	//}
	//
	//res, err = proto.Marshal(&open_shifts)

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	return shifts
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func getReplacements(manager Manager) []Employee {
	myrand := random(1, 6)
	replacements := []Employee{}
	for i := 0; i < myrand; i++ {
		e := generateEmmployee()

		replacements = append(replacements, e)
	}
	return replacements
}
func generateEmmployee() Employee {
	e := Employee{
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
	}

	return e
}
