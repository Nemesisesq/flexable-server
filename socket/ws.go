package socket

import (
	"encoding/json"
	"net/http"
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

	ws *websocket.Conn
}

func NewHandler() *Handler {
	return &Handler{
		out: make(chan []byte),
	}
}

func (h Handler) recieveData(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	log.Info("Connecting to recieve data")
	for {

		// read messages

		_, message, err := ws.ReadMessage()
		log.Info("got a message")
		log.Info(string(message))
		if err != nil {
			break
		}

		// Do stuff

		go func() {

			message, err = json.Marshal(time.Now())
			if err != nil {
				log.Error(err)
			}

			h.out <- message
		}()

		// Write message out to app

	}
}

func (h Handler) sendData(ws *websocket.Conn, done chan struct{}) {
	defer func() {
	}()

	// blocking until error
	for {
		ws.SetWriteDeadline(time.Now().Add(writeWait))

		select {
		case m := <-h.out:
			if err := ws.WriteMessage(websocket.TextMessage, m); err != nil {
				h.internalError("something has gone wrong", err)
				close(h.out)
				ws.Close()
				break
			}

		}
	}

	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
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

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	h.ws = ws
	defer ws.Close()

	stdoutDone := make(chan struct{})
	go h.sendData(ws, stdoutDone)
	go ping(ws, stdoutDone)

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

}
