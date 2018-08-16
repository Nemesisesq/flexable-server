package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/nemesisesq/flexable-server/account"
	"github.com/nemesisesq/flexable-server/db"
	"github.com/nemesisesq/flexable-server/flexable"
	"github.com/nemesisesq/flexable-server/shifts"
	"github.com/odknt/go-socket.io"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"github.com/x-cray/logrus-prefixed-formatter"
	"net/http/httputil"
	"context"
)

func main() {
	log.SetFormatter(new(prefixed.TextFormatter))
	log.SetLevel(log.DebugLevel)
	port := os.Getenv("PORT")
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize global mongo database
	session, _, _ := db.InitDB()

	flexable.SocketServerConnections(*server, "flexable")
	//flexable.SocketServerConnections(*server, "employee")
	flexable.SetListeners(server)

	go server.Serve()
	defer server.Close()

	m := mux.NewRouter()
	r := render.New()
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

	n.UseHandler(m)

	//Mongo database
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {

		ctx := r.Context()
		sesh := session.Clone()
		defer sesh.Close()
		ctx = context.WithValue(ctx, "db", sesh)
		r = r.WithContext(ctx)
		n(w, r)
	})

	m.Handle("/socket.io/", server)

	m.HandleFunc(fmt.Sprintf("/profile/%d", flexable.SET_PROFILE),
		func(writer http.ResponseWriter, request *http.Request) {

	})

	selectVolunteerEndpoint := fmt.Sprintf("/manager/%d", flexable.SELECT_VOLUNTEER)
	m.HandleFunc(selectVolunteerEndpoint,
		func(writer http.ResponseWriter, request *http.Request) {

			shifts.SelectVolunteer(request)
			fmt.Println("Chosing Volunteer")

			b, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}

			println(string(b))
		})

	// register http handlers
	for _, v := range flexable.HttpTypes {
		m.HandleFunc(fmt.Sprintf("/%v/%d", v.N, v.T), v.H)
	}

	m.HandleFunc("/users/push-token", func(writer http.ResponseWriter, request *http.Request) {
		err := account.SavePushToken(*request)
		if err != nil {
			r.JSON(writer, http.StatusInternalServerError, map[string]interface{}{"error": err})
		}

		r.JSON(writer, http.StatusOK, map[string]interface{}{"hello": "world"})
	})

	m.HandleFunc("/users/verify", func(writer http.ResponseWriter, request *http.Request) {

		log.Info("registerng user")
		user := account.UserRole(*request)

		r.JSON(writer, http.StatusOK, map[string]interface{}{"role": user.Role, "profile": user.Profile, "is_admin": user.IsAdmin})
		return
	})

	m.HandleFunc("/sms/incoming/{smsId}", func(writer http.ResponseWriter, request *http.Request) {
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
	m.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, n))
}
