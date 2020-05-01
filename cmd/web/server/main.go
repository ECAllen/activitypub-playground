package main

import (
	"crypto/rsa"
	"fmt"
	pub "github.com/go-ap/activitypub"
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/middleware"
	"crypto/sha1"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Globals
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

type application struct {
	people []*pub.Person
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// TODO use context or dependency injection?
	/* ctx := r.Context()
	article, ok := ctx.Value("article").(*Article)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	*/
	w.Write([]byte("OK"))
}

/*
	Person {
		  ID: http://127.0.0.1:4000/actors/35318264c9a98faf79965c270ac80c5606774df1
		  Type:Person
		  Name:Alice
		  Attachment:<nil>
		  AttributedTo:<nil>
		  Audience:[]
		  Content:[]
		  Context:<nil>
		  MediaType:
		  EndTime:0001-01-01 00:00:00 +0000 UTC
		  Generator:<nil>
		  Icon:<nil>
		  Image:<nil>
		  InReplyTo:<nil>
		  Location:<nil>
		  Preview:<nil>
		  Published:0001-01-01 00:00:00 +0000 UTC
		  Replies:<nil>
		  StartTime:0001-01-01 00:00:00 +0000 UTC
		  Summary:[]
		  Tag:[]
		  Updated:0001-01-01 00:00:00 +0000 UTC
		  URL:<nil>
		  To:[]
		  Bto:[]
		  CC:[]
		  BCC:[]
		  Duration:0s
		  Likes:<nil>
		  Shares:<nil>
		  Source:{Content:[] MediaType:}
		  Inbox:0xc0000ce580
		  Outbox:0xc0000ce840
		  Following:<nil>
		  Followers:<nil>
		  Liked:0xc0000cf080
		  PreferredUsername:[]
		  Endpoints:<nil>
		  Streams:[]
		  PublicKey:{ID: Owner:<nil> PublicKeyPem:}
		}
*/

func (app *application) actors(w http.ResponseWriter, r *http.Request) {
	for _, obj := range app.people {
		fmt.Printf("%+v\n", obj)
	}
	ppl, err := json.MarshalIndent(app.people, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ppl)
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.home)
	// TODO double check standard for endpoints
	r.Get("/actors", app.actors)
	return r
}

func main() {

	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	var ppl []*pub.Person
	app := &application{
		people: ppl,
	}

	names := []string{"Alice", "Bob"}
	for _, name := range names {
		actorName := pub.NaturalLanguageValues{{Ref: pub.NilLangRef, Value: name}}
		h := sha1.New()
		h.Write([]byte(name))
		actorSumHex := h.Sum(nil)
		actorSum := fmt.Sprintf("%x", actorSumHex)
		actorsURL := fmt.Sprintf("%s/actors", apiURL)
		actorURL := fmt.Sprintf("%s/%s", actorsURL, actorSum)

		actorInbox := fmt.Sprintf("%s/inbox", actorURL)
		in := pub.OrderedCollectionNew(pub.ID(actorInbox))

		actorOutbox := fmt.Sprintf("%s/outbox", actorURL)
		out := pub.OrderedCollectionNew(pub.ID(actorOutbox))

		id := pub.IRI(actorURL)
		a := pub.PersonNew(id)
		a.Name = actorName
		a.Inbox = in
		a.Outbox = out
		app.people = append(app.people, a)
	}

	srv := &http.Server{
		Addr:         host,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	logger.WithFields(logrus.Fields{
		"event": "ListenAndServeTLS",
	}).Fatal(err)
}

/* TODO

# Server

receive message and put in inbox

# Client

send message

public key
test by sending message to server
*/
