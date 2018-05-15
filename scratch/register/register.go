package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func PostRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hoi")
}

func GetRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hoi")
}

type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	ProjectID    string
}

var (
	config    configuration
	appRouter mux.Router
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LoadConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
}

/*
Standard http handler
*/
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

/*
Our context.Context http handler
*/
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

/*
  Higher order function for changing a HandlerFunc to a ContextHandlerFunc,
  usually creating the context.Context along the way.
*/
type ContextHandlerToHandlerHOF func(f ContextHandlerFunc) HandlerFunc

/*
Creates a new Context and uses it when calling f
*/
func ContextHanderToHttpHandler(f ContextHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		f(ctx, w, r)
	}
}

/*
Creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFuncs.
*/
func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/register", f(PostRegistrationHandler)).Methods("POST")
	appRouter.HandleFunc("/register", f(GetRegistrationHandler)).Methods("GET")

	return appRouter
}

func init() {
	LoadConfig()
	http.Handle("/", CreateHandler(ContextHanderToHttpHandler))
}
