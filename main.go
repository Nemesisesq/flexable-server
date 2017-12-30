package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/nemesisesq/flexable/flexable"
	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {

	port := os.Getenv("PORT")
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	flexable.SocketServerConnections(*server)

	flexable.SetListeners(server)
	go server.Serve()
	defer server.Close()

	r := mux.NewRouter()
	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)
	r.Handle("/socket.io/", server)
	r.HandleFunc("/sms/listener/{smsId}", func(writer http.ResponseWriter, request *http.Request) {
		log.Info("Message response received!")
	})
	r.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":"+port, n))
}
