package flexable

import (
	"context"
	"time"

	"github.com/mitchellh/hashstructure"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	//log "github.com/sirupsen/logrus"
	"github.com/globalsign/mgo/bson"
)

func InitWatchers(socket socketio.Conn) {

	go func() {
		CheckOpenShifts(socket)
	}()

}
func CheckOpenShifts(s socketio.Conn) {

	ctx := s.Context().(context.Context)

	ticker := time.NewTicker(time.Second * 2)
	tickChan := ticker.C
	companyId := "123"
	var currentShiftState uint64

	for {
		shiftList := []shifts.Shift{}
		select {
		case <-tickChan:
			//log.Info("I'm checking the db")
			//log.Info(db)
			shiftList = shifts.GetAllShifts(bson.M{"company.uuid": companyId})

			shift_list_hash, err := hashstructure.Hash(&shiftList, nil)

			if err != nil {
				panic(err)
			}

			if currentShiftState != shift_list_hash {
				currentShiftState = shift_list_hash
				s.Emit(constructSocketID(OPEN_SHIFTS), shiftList)
			}

		case <-ctx.Done():
			ticker.Stop()
			break
		}

	}

}
