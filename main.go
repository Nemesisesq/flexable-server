package main

import (
	"flag"

	"github.com/desertbit/glue"
	log "github.com/sirupsen/logrus"

	_ "github.com/heroku/x/hmetrics/onload"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// Create a new glue server.
	server := glue.NewServer(glue.Options{
		HTTPListenAddress: ":8080",
	})

	// Release the glue server on defer.
	// This will block new incoming connections
	// and close all current active sockets.
	defer server.Release()

	// Set the glue event function to handle new incoming socket connections.
	server.OnNewSocket(onNewSocket)

	//socket.OpenShifts()

	// Run the glue server.

	log.Info("Running Glue server")
	err := server.Run()
	if err != nil {
		log.Fatalf("Glue Run: %v", err)
	}

}
func OpenShiftsChan(s *glue.Socket) {
	c := s.Channel("golang")
	// Set the channel on read event function.
	c.OnRead(func(data string) {
		// ...
		log.Info("new message from golang channel")
	})
	// Write to the channel.
	c.Write("Hello Gophers!")
}

func onNewSocket(s *glue.Socket) {
	// Set a function which is triggered as soon as the socket is closed.
	s.OnClose(func() {
		log.Printf("socket closed with remote address: %s", s.RemoteAddr())
	})

	OpenShiftsChan(s)

	// Set a function which is triggered during each received message.
	s.OnRead(func(data string) {
		// Echo the received data back to the client.
		s.Write(data)
	})

	// Send a welcome string to the client.
	s.Write("Hello Client")
}
