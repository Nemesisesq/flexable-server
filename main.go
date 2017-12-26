package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/nemesisesq/flexable/employees"
	"github.com/nemesisesq/flexable/protobuf"
	"github.com/nemesisesq/flexable/socket"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/nats-io/go-nats"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	r := mux.NewRouter()

	n := negroni.Classic()

	r.HandleFunc("/", home)

	n.Use(c)
	n.UseHandler(r)

	NatsListen()
	log.Print("Listening on port ", *addr)
	log.Fatal(http.ListenAndServe(*addr, n))

}
func NatsListen() {
	log.Info("star" +
		"ting nats listener")
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		panic(err)
	}

	chan0 := fmt.Sprintf("flexable.data.service.%v", int(payload.Payload_OPEN_SHIFTS))
	nc.QueueSubscribe(chan0, FormQueue(payload.Payload_OPEN_SHIFTS), func(m *nats.Msg) {
		log.Info("got request for open shifts")
		p := &payload.Payload{}
		shifts := socket.OpenShifts(p)
		log.Info("send response for open shifts")

		out, err := json.Marshal(&shifts)
		if err != nil {
			panic(err)
		}
		nc.Publish(m.Reply, out)
	})

	nc.QueueSubscribe(fmt.Sprintf("flexable.communication.service.%v", int(payload.Payload_FIND_SHIFT_SUBSTITUTE)), FormQueue(payload.Payload_FIND_SHIFT_SUBSTITUTE), func(m *nats.Msg) {

		log.Info("Got a request to find shift substitute")

		// Get Available Employees

		shiftDetails := map[string]interface{}{}

		err := json.Unmarshal(m.Data, &shiftDetails)

		if err != nil {
			panic(err)
		}

		emps := employees.Employees{}

		emps.GetAvailable()

		for _, v := range emps {
			v.Contact(shiftDetails)
		}

	})

}

func FormQueue(pt payload.Payload_Type) string {
	return strconv.Itoa(int(pt))
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/stream")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
