package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	//"github.com/auth0/go-jwt-middleware"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/nemesisesq/flexable/account"
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

	/*jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})*/

	n.UseHandler(r)
	r.Handle("/socket.io/", server)

	r.HandleFunc("/users/push-token", func(writer http.ResponseWriter, request *http.Request) {
		account.SavePushToken(*request)
	})

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

		if res != nil {

			out, err := xml.MarshalIndent(res, "", "   ")

			if err != nil {
				panic(err)
			}

			fmt.Fprint(writer, out)
		}

	})
	r.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, n))
}
