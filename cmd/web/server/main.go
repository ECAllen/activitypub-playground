package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	pub "github.com/go-ap/activitypub"
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/middleware"
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
	/* ctx := r.Context()
	article, ok := ctx.Value("article").(*Article)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}
	*/
	w.Write([]byte("OK"))
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", app.home)
	return r
}
func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger.Formatter = new(logrus.JSONFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	var ppl []*pub.Person
	app := &application{
		people: ppl,
	}

	srv := &http.Server{
		Addr:         *addr,
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
-  create server
-  create endpoints
-  listen on endpoints

respond with list of actors to client
receive message and put in inbox

# Client

- create client
client gets list of actors
client get inbox endpoint of actor
publickey
send message

test by sending message to server
*/
