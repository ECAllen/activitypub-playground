package main

import (
	"crypto/rsa"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	pub "github.com/go-ap/activitypub"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
)

var (
	UserAgent    = "ap-go-http-client"
	HeaderAccept = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`

	host            = "127.0.0.1:4000"
	apiURL          = "http://127.0.0.1:4000"
	authCallbackURL = fmt.Sprintf("%s/auth/local/callback", apiURL)

	rnd    = rand.New(rand.NewSource(6667))
	key, _ = rsa.GenerateKey(rnd, 512)

	logger = logrus.New()
)

func inbox(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{"event": "Inbox"}).Trace("Inbox hit")
}

func actors(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{"event": "Actor"}).Trace("Actor hit")
	// TODO

	/*
		get list of actors
		json marshall actors
		respond with json
		items: map[string]*objectVal{
									"e869bdca-dd5e-4de7-9c5d-37845eccc6a1": {
										id:      "http://127.0.0.1:9998/actors/e869bdca-dd5e-4de7-9c5d-37845eccc6a1",
										typ:     string(pub.PersonType),
										summary: "Generated actor",
										content: "Generated actor",
										url:     "http://127.0.0.1:9998/actors/e869bdca-dd5e-4de7-9c5d-37845eccc6a1",
										inbox: &objectVal{
											id: "http://127.0.0.1:9998/actors/e869bdca-dd5e-4de7-9c5d-37845eccc6a1/inbox",
										},
										outbox: &objectVal{
											id: "http://127.0.0.1:9998/actors/e869bdca-dd5e-4de7-9c5d-37845eccc6a1/outbox",
										},
										liked: &objectVal{
											id: "http://127.0.0.1:9998/actors/e869bdca-dd5e-4de7-9c5d-37845eccc6a1/liked",
										},
										preferredUsername: "johndoe",
										name:              "Johnathan Doe",
									},
								},
	*/

}

func handlerReqs() {
	http.HandleFunc("/actors", actors)
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

type Actors struct {
	Total  int
	People []*pub.Person
}

func (a *Actors) Add(person *pub.Person) {
	a.Total++
	a.People = append(a.People, person)
	return
}

func main() {
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	actorName := "Alice"
	h := sha1.New()
	h.Write([]byte(actorName))
	actorSumHex := h.Sum(nil)
	actorSum := fmt.Sprintf("%x", actorSumHex)

	actorsURL := fmt.Sprintf("%s/actors", apiURL)
	actorURL := fmt.Sprintf("%s/%s", actorsURL, actorSum)

	logger.WithFields(logrus.Fields{
		"event":      "SHA1 Sum",
		"Actor Name": actorName,
		"Actor Sum":  actorSum,
		"Actor URL":  actorURL,
	}).Trace("Actor Sum")

	actors := Actors{}
	id := pub.IRI(actorURL)
	a := pub.PersonNew(id)
	actors.Add(a)

	aj, _ := json.MarshalIndent(a, "", "  ")
	logger.WithFields(logrus.Fields{
		"event": "Create",
		"actor": string(aj),
	}).Trace("New Person")

	handlerReqs()
}

/* TODO

# Server
-  create server
-  create endpoints
-  listen on endpoints

respond with list of Actors to client
receive message and put in inbox

# Client

- create client
client gets list of actors
client get inbox endpoint of actor
publickey
send message

test by sending message to server
*/
