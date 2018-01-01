package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/nemesisesq/flexable/flexable"
	"github.com/nemesisesq/flexable/shifts"
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
	r.HandleFunc("/sms/incoming/{smsId}", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request)
		vars := mux.Vars(request)
		smsID := vars["smsId"]
		body := map[string]string{}
		request.ParseForm()
		for key, values := range request.Form { // range over map
			for _, value := range values { // range over []string
				body[key] = value
			}
		}
		res := shifts.UpdateShiftWithSmsID(smsID, body)

		if res {

			fmt.Fprint(writer, "Thanks for voluteering we will be getting back to you shortly to let you know you got the gig!")
		} else {
			fmt.Fprint(writer, "Sorry, this shift wasn't for you we'll chdeheck with other opportunities")
		}

	})
	r.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, n))
}
