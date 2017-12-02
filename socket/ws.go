package socket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

type Handler struct {
	//messages to the client
	out chan []byte
	mu  sync.Mutex
	ws  *websocket.Conn
}

func NewHandler() *Handler {

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in f", r)
		}
	}()

	return &Handler{
		out: make(chan []byte),
	}
}

func (h Handler) recieveData(ws *websocket.Conn) {

	//defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	err := ws.SetReadDeadline(time.Now().Add(pongWait))

	if err != nil {
		panic(err)
	}
	ws.SetPongHandler(func(string) error {
		err := addPongWait(ws)

		if err != nil {
			return err
		}
		return nil
	})

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	//log.Info("Connecting to recieve data")
	for {

		// read messages

		_, message, err := ws.ReadMessage()

		if err != nil {
			break
		}

		if string(message) == "__ping__" {
			log.Info("__pong__")
			h.out <- []byte("__pong__")
		} else if len(message) != 0 {

			log.Info("got a message", string(message))
			log.Info(string(message))

			// Do stuff

			func() {
				log.Info("Handling message")

				//res := HandleMessage(message)
				log.Info("I get here")
				//h.out <- res
			}()

			// Write message out to app
		}

	}

	log.Info("exiting recieveData")
}
func addPongWait(ws *websocket.Conn) error {
	return ws.SetReadDeadline(time.Now().Add(pongWait))
}

func (h *Handler) send(v interface{}) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.ws.WriteJSON(v)
}

func (h Handler) sendData(ws *websocket.Conn, done chan struct{}) {

	// blocking until error
	for {
		ws.SetWriteDeadline(time.Now().Add(writeWait))

		select {
		case m := <-h.out:
			log.Info("writing message from channel")

			if err := h.send(m); err != nil {
				h.internalError("something has gone wrong", err)
				//close(h.out)
			}

		}
	}

	//close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	//ws.Close()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
}

func (h Handler) internalError(msg string, err error) {
	log.Println(msg, err)
	h.ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h Handler) ServeWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)

	log.Info("connecting")

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	h.ws = ws
	//defer ws.Close()

	stdoutDone := make(chan struct{})
	go h.sendData(ws, stdoutDone)

	// blocking with receive data
	h.recieveData(ws)
	log.Info("exiting")
	select {
	case <-stdoutDone:
	case <-time.After(time.Second):
		// A bigger bonk on the head.
		//if err := proc.Signal(os.Kill); err != nil {
		//	log.Println("term:", err)
		//}
		<-stdoutDone
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
}
