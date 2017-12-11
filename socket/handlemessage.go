package socket

import (
	"math/rand"
	"time"

	"github.com/manveru/faker"
	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/employees"
	"github.com/nemesisesq/flexable/manager"
	payload2 "github.com/nemesisesq/flexable/protobuf"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/nemesisesq/flexable/user"
	log "github.com/sirupsen/logrus"
)

type Location interface{}

var fake, _ = faker.New("en")

func OpenShifts(p *payload2.Payload) (res interface{}) {

	mgr := manager.Manager{
		user.User{
			Name: user.Name{
				FirstName: fake.FirstName(),
				LastName:  fake.LastName(),
				NickNames: []string{fake.JobTitle()},
			},
			Email: fake.Email(),
			PhoneNumber: user.PhoneNumber{
				Temp: fake.PhoneNumber(),
			},
		},
	}

	comp := company.Company{
		fake.CompanyName(),
		fake.StreetAddress(),
	}

	var sh []shifts.Shift
	for i := 0; i < 4; i++ {

		employee := generateEmmployee()
		s := shifts.Shift{
			Employee: employee,
			Manager:  mgr,
			ShiftSite: shifts.ShiftSite{
				fake.JobTitle(),
				fake.StreetAddress(),
				comp,
				[]manager.Manager{mgr},
			},
			DateTime:     time.Now().Add(time.Hour * 2),
			Replacements: getReplacements(mgr),
		}
		sh = append(sh, s)

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

	return sh
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func getReplacements(manager manager.Manager) []employees.Employee {
	myrand := random(1, 6)
	replacements := []employees.Employee{}
	for i := 0; i < myrand; i++ {
		e := generateEmmployee()

		replacements = append(replacements, e)
	}
	return replacements
}
func generateEmmployee() employees.Employee {
	e := employees.Employee{
		User: user.User{
			user.Name{
				fake.FirstName(),
				fake.LastName(),
				[]string{fake.JobTitle()},
			},
			fake.Email(),
			user.PhoneNumber{
				Temp: fake.PhoneNumber(),
			},
		},
	}

	return e
}
