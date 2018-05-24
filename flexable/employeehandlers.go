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
	"fmt"
)

type EmployeeData struct {
}

func PickUpShift(s socketio.Conn, data interface{}) {

	log.Info("Picking up  a shift")
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

	//	Notify Manager

	template := `Hey There is an new volunteer for the shift from {{.StartTime }} to {{.EndTime}} On {{.Date }}!`

	buf, err := CreateTextMessageString(template, shift)

	if err != nil {
		panic(err)
	}

	flexableAdmin := account.User{Email:"admin@myflexable.com", Role:"admin"}

	shift.Manager.Notify([]string{buf.String(),buf.String()}, "Someone Volunteered for the shift you created!", shift.PhoneNumber, flexableAdmin)

}

func UpdateProfile(s socketio.Conn, data interface{}) {
	log.Info("Updating Profile")
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

}

func CallOfShift(s socketio.Conn, data interface{}) {

	log.Info("CAlling off shift")
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
		if v.Profile.Email == emp.Profile.Email {
			shift.Volunteers = append(shift.Volunteers[:k], shift.Volunteers[k+1:]...)
		}
	}

	if shift.Chosen.Profile.Email == emp.Profile.Email {
		shift.Chosen = account.User{}
	}
	shift.Save()

	// Notify manager that the person chosen for this shift has called off
	template := fmt.Sprintf(`Uh Oh the shift from {{.StartTime }} to {{.EndTime}} On {{.Date }} has been called off by %v %v (%v)`, user.Profile.FirstName, user.Profile.LastName, user.Profile.Email)


	buf, err := CreateTextMessageString(template, shift)

	if err != nil {
		panic(err)
	}

	flexableAdmin := account.User{Email:"admin@myflexable.com", Role:"admin"}

	shift.Manager.Notify([]string{buf.String(),buf.String()}, "Someone Called Off!", shift.PhoneNumber, flexableAdmin)


}

func GetOpenShifts(s socketio.Conn, data interface{}) {
	log.Info("Returning  employee openshifts")

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User);
	ticker := time.NewTicker(time.Second)
	timeout := time.NewTicker(time.Minute)

	tickerChan := ticker.C
	var currentShiftState uint64
	//<-tickerChan
	go func() {
		for {
			shiftList := []shifts.Shift{}
			select {
			case <-tickerChan:
				var query bson.M
				if user.Profile.Company.UUID == "" {
					user.Profile.Company.UUID = "123"
				}
				query = bson.M{"company.uuid": user.Profile.Company.UUID}
				shiftList = shifts.GetAllShifts(query)
				shift_list_hash, err := hashstructure.Hash(&shiftList, nil)

				if err != nil {
					panic(err)
				}
				//println(shift_list_hash)
				if currentShiftState != shift_list_hash {
					log.Info("Employee shifts are updating ")
					currentShiftState = shift_list_hash
					s.Emit(constructSocketID(GET_OPEN_SHIFTS), shiftList, func(so socketio.Conn, data string) {
						log.Println("Client ACK with data: ", data)
					})
				}

			case <-ctx.Done():
				ticker.Stop()
				return
			case <-timeout.C:
				log.Info("I'm closing out the channel")
				cancel := ctx.Value("cancel").(context.CancelFunc)
				cancel()
				s.Close()
			}
		}
	}()
	log.Info("Exiting go loop")
}

func UpdateNotifications(s socketio.Conn, data interface{}) {
	log.Info("Broadcasting Notifications")

	ctx := s.Context().(context.Context)
	user := ctx.Value("user").(account.User);
	ticker := time.NewTicker(time.Second)
	timeout := time.NewTimer(time.Minute)
	var currentNotificationState uint64
	//<-tickerChan
	go func() {
	L:
		for {
			select {
			case <-ticker.C:
				var query bson.M
				query = bson.M{"_id": user.ID}
				out := user.Find(query)

				notifications := out.Notifications
				notifications_hash, err := hashstructure.Hash(&notifications, nil)

				if err != nil {
					panic(err)
				}
				//println(shift_list_hash)
				if currentNotificationState != notifications_hash {
					log.Info("notifications are updating ")
					currentNotificationState = notifications_hash
					s.Emit(constructSocketID(UPDATE_NOTIFICATIONS), notifications, func(so socketio.Conn, data string) {
						log.Println("Client ACK with data: ", data)

					})
					timeout = time.NewTimer(time.Minute)
				}

			case <-ctx.Done():
				ticker.Stop()
				log.Info("I'm stopping the ticker now")
				break L

			case <-timeout.C:
				log.Info("I'm closing out the channel")
				cancel := ctx.Value("cancel").(context.CancelFunc)
				cancel()
				s.Close()
			}
		}
		log.Info("Exiting go loop")
	}()
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
