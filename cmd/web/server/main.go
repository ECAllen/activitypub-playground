package main

import (
	"encoding/json"
	"github.com/go-ap/activitypub"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func inbox(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{"event": "Inbox"}).Trace("Inbox hit")
}

func actor(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{"event": "Actor"}).Trace("Actor hit")
}

func handlerReqs() {
	http.HandleFunc("/actor", actor)
	http.HandleFunc("/inbox", inbox)
	logger.WithFields(logrus.Fields{"event": "Handle Requests"}).Trace("Handle Reqs")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "Err Handle Requests",
			"error": err,
		}).Fatal("Handle Reqs")
	}
}

var logger = logrus.New()

func main() {
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	id := activitypub.IRI("http://localhost:4000/actor")
	a := activitypub.PersonNew(id)
	aj, _ := json.MarshalIndent(a, "", "  ")

	logger.WithFields(logrus.Fields{
		"event": "Create",
		"actor": string(aj),
	}).Trace("New Person")

	handlerReqs()
}

/* TODO

- create server
-  create endpoints
-  listen on endpoints
key handling???
  receive message and put in inbox

create client
  publickey
  send message

test by sending message to server
*/
