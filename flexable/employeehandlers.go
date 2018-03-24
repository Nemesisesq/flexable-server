package flexable

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/hashstructure"
	"github.com/nemesisesq/flexable/employee"
	payload2 "github.com/nemesisesq/flexable/protobuf"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
)

type EmployeeData struct {
}

func PickUpShift(s socketio.Conn, data interface{}) interface{} {
	return nil
}
func CallOfShift(s socketio.Conn, data interface{}) interface{} {
	return nil
}

func GetOpenShifts(s socketio.Conn, data interface{}) interface{} {
	return nil
}
func GetEmployeeShifts(s socketio.Conn, data interface{}) interface{} {

	payload := data.(map[string]interface{})["payload"]

	if payload == nil {
		return nil
	}
	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	var empl_payload payload2.EmployeePayload
	err = json.Unmarshal(tmp, &empl_payload)

	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(time.Second * 2)
	tickChan := ticker.C
	//companyId := "123"
	var employeeShiftHash uint64

	go func() {
		log.Info("starting employee shift watcher")

		for {
			select {
			case <-tickChan:

				shift_list := shifts.GetAllShifts(bson.M{
					"$and": []bson.M{
						{"volunteers": bson.M{"$size": 0}},
						{"company.uuid": 123},
					}})

				currentEmployee := employee.GetOneEmployee(bson.M{"id": empl_payload.Id})

				//Combine the employee shift list as well as make them unique

				var shift_list2 []employee.Shiftable

				// Loop through and cast shift.Shift to Shiftable
				for _, v := range shift_list {
					var tmp employee.Shiftable = v
					shift_list2 = append(shift_list2, tmp)
				}
				combined := append(shift_list2, currentEmployee.Schedule...)

				unique := []shifts.Shift{}

				for _, v := range combined {
					found := false
					for _, uv := range unique {
						if v.(shifts.Shift).Date == uv.Date {
							found = true
						}
					}

					if !found {
						unique = append(unique, v.(shifts.Shift))
					}

				}

				tmpHash, err := hashstructure.Hash(unique, nil)

				if err != nil {
					panic(err)
				}

				if tmpHash != employeeShiftHash {
					s.Emit(constructSocketID(GET_MY_SHIFTS), unique)
				}

			}

		}

	}()
	return nil
}
