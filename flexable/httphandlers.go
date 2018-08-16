package flexable

import (
	"net/http"
	"github.com/manveru/faker"
	"github.com/nemesisesq/flexable-server/position"
	"github.com/nemesisesq/flexable-server/utils"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func GetPositionsHttp(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting positions")
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	jobs := position.GetAllPositions(nil)

	if len(jobs) <= 1 {

		for i := 0; i < utils.RandomRange(2, 10); i++ {
			x := position.Position{
				bson.NewObjectId(),
				fake.JobTitle(),
				10.00,
				"hr",
			}

			x.Save()
			jobs = append(jobs, x)
		}
	}

	fmt.Println(jobs)

	encoder := json.NewEncoder(w)
	encoder.Encode(jobs)

}
