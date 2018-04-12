package flexable

import (
	"encoding/json"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
	"context"
	"github.com/nemesisesq/flexable/account"
	"github.com/mitchellh/hashstructure"
)

type EmployeeData struct {
}

func PickUpShift(s socketio.Conn, data interface{}) interface{} {

	log.Info("Picking up  a shift" )
	payload := data.(map[string]interface{})["payload"]

	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User)

	var shift shifts.Shift
	err = json.Unmarshal(tmp, &shift)

	if err != nil {
		panic(err)
	}

	emp := user


	shift.Volunteers = append(shift.Volunteers, emp)

	shift.Save()



	return nil
}

func UpdateProfile(s socketio.Conn, data interface{}) interface{} {
	log.Info("CAlling off shift" )
	payload := data.(map[string]interface{})["payload"]

	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User)

	err = json.Unmarshal(tmp, &user.Profile)

	if err != nil {
		panic(err)
	}

	user.Upsert(bson.M{"_id": user.ID})

	s.Emit(constructSocketID(SET_PROFILE), &user.Profile)


	return nil
}


func CallOfShift(s socketio.Conn, data interface{}) interface{} {


	log.Info("CAlling off shift" )
	payload := data.(map[string]interface{})["payload"]

	tmp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User)

	var shift shifts.Shift
	err = json.Unmarshal(tmp, &shift)

	if err != nil {
		panic(err)
	}

	emp := user

	for k, v := range shift.Volunteers {
		if v.Profile.Email ==  emp.Profile.Email {
			shift.Volunteers = append(shift.Volunteers[:k], shift.Volunteers[k+1:]...)
		}
	}

	if shift.Chosen.Profile.Email == emp.Profile.Email{
		shift.Chosen = account.User{}
	}
	shift.Save()
	return nil
}

func GetOpenShifts(s socketio.Conn, data interface{}) interface{} {
	log.Info("Returning  employee openshifts")

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User);
	ticker := time.NewTicker(time.Second)

	tickerChan := ticker.C
	var currentShiftState uint64
	//<-tickerChan
	for {
		shiftList := []shifts.Shift{}
		select {
		case <-tickerChan:
			var query bson.M
			query = bson.M{"company.uuid": user.Profile.Company.UUID}
			shiftList = shifts.GetAllShifts(query)
			shift_list_hash, err := hashstructure.Hash(&shiftList, nil)

			if err != nil {
				panic(err)
			}
			if currentShiftState != shift_list_hash {
				currentShiftState = shift_list_hash
				s.Emit(constructSocketID(GET_OPEN_SHIFTS), shiftList, func(so socketio.Conn, data string) {
					log.Println("Client ACK with data: ", data)
				})
			}
		}
	}
	return "hello"
}

//func GetEmployeeShifts(s socketio.Conn, data interface{}) interface{} {
//	log.Info("getting employee shifts")
//	payload := data.(map[string]interface{})["payload"]
//
//	if payload == nil {
//		return nil
//	}
//	tmp, err := json.Marshal(payload)
//	if err != nil {
//		panic(err)
//	}
//	var empl_payload payload2.EmployeePayload
//	err = json.Unmarshal(tmp, &empl_payload)
//
//	if err != nil {
//		panic(err)
//	}
//	ticker := time.NewTicker(time.Second * 2)
//	tickChan := ticker.C
//	//companyId := "123"
//	var employeeShiftHash uint64
//
//	//<-tickChan
//
//	go func() {
//		log.Info("starting employee shift watcher")
//
//		for {
//			select {
//			case <-tickChan:
//
//				shift_list := shifts.GetAllShifts(bson.M{
//					"$and": []bson.M{
//						{"volunteers": bson.M{"$size": 0}},
//						{"company.uuid": 123},
//					}})
//
//				currentEmployee := employee.GetOneEmployee(bson.M{"id": empl_payload.Id})
//
//				//Combine the employee shift list as well as make them unique
//
//				var shift_list2 []employee.Shiftable
//
//				// Loop through and cast shift.Shift to Shiftable
//				for _, v := range shift_list {
//					var tmp employee.Shiftable = v
//					shift_list2 = append(shift_list2, tmp)
//				}
//				combined := append(shift_list2, currentEmployee.Schedule...)
//
//				unique := []shifts.Shift{}
//
//				for _, v := range combined {
//					found := false
//					for _, uv := range unique {
//						if v.(shifts.Shift).Date == uv.Date {
//							found = true
//						}
//					}
//
//					if !found {
//						unique = append(unique, v.(shifts.Shift))
//					}
//
//				}
//
//				tmpHash, err := hashstructure.Hash(unique, nil)
//
//				if err != nil {
//					panic(err)
//				}
//
//				if tmpHash != employeeShiftHash {
//					s.Emit(constructSocketID(GET_MY_SHIFTS), unique)
//				}
//
//			}
//
//		}
//
//	}()
//	return nil
//}
