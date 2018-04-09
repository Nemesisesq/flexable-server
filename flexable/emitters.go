package flexable

import (
	"context"
	"time"

	"github.com/mitchellh/hashstructure"
	"github.com/nemesisesq/flexable/shifts"
	"github.com/odknt/go-socket.io"
	//log "github.com/sirupsen/logrus"
	"github.com/globalsign/mgo/bson"
	"github.com/oxequa/grace"
	"github.com/nemesisesq/flexable/account"
	log"github.com/sirupsen/logrus"
)

func InitWatchers(socket socketio.Conn) {

	go func() {
		CheckOpenShifts(socket)
		TestPushNotifications(socket)
	}()

}

func TestPushNotifications(s socketio.Conn) {
	ctx := s.Context().(context.Context)

	user := ctx.Value("user").(account.User);
	ticker := time.NewTicker(time.Second * 90)

	tickerChan := ticker.C

	for {
		select {
		case <-tickerChan:
			log.Info("sending push message")
			user.Notify("This is a test message welcome to the family", "Test title", &shifts.Shift{})
		}
	}
}


func CheckOpenShifts(s socketio.Conn)(e error) {
defer grace.Recover(&e)
	ctx := s.Context().(context.Context)

	ticker := time.NewTicker(time.Second * 2)
	tickChan := ticker.C
	companyId := "123"
	var currentShiftState uint64

	for {
		shiftList := []shifts.Shift{}
		select {
		case <-tickChan:
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
