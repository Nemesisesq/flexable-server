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
	log "github.com/sirupsen/logrus"
)

func InitWatchers(socket socketio.Conn) {

	go func() {

		go CheckOpenShifts(socket)

		//TestPushNotifications(socket)
	}()

}

func TestPushNotifications(s socketio.Conn) {
	log.Debug("I'm testing push notifications on an interval")
	ctx := s.Context().(context.Context)

	user := ctx.Value("user").(account.User);
	ticker := time.NewTicker(time.Hour * 2)

	tickerChan := ticker.C

	for {
		select {
		case <-tickerChan:

			log.Debug("Firing off!! ")
			log.Info("sending push message")
			user.Notify("This is a test message welcome to the family", "Test title", shifts.Shift{}.PhoneNumber)
		}
	}
}

func CheckOpenShifts(s socketio.Conn) (e error) {
	defer grace.Recover(&e)
	ctx := s.Context().(context.Context)

	ticker := time.NewTicker(time.Second * 2)
	timeout := time.NewTimer(time.Minute)
	companyId := "123"
	var currentShiftState uint64
L:
	for {
		shiftList := []shifts.Shift{}
		select {
		case <-ticker.C:
			shiftList = shifts.GetAllShifts(bson.M{"company.uuid": companyId})
			shift_list_hash, err := hashstructure.Hash(&shiftList, nil)

			if err != nil {
				panic(err)
			}

			if currentShiftState != shift_list_hash {
				currentShiftState = shift_list_hash
				s.Emit(constructSocketID(OPEN_SHIFTS), shiftList)
				timeout = time.NewTimer(time.Minute)
			}

		case <-ctx.Done():
			ticker.Stop()
			log.Info("I'm stopping the ticker")
			break L

		case <-timeout.C:
			log.Info("I'm closing out the channel")
			s.Close()
			cancel := ctx.Value("cancel").(context.CancelFunc)
			cancel()
		}
	}
	log.Info("exiting go routine")
	return e
}
