package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

type Registration struct {
	First_Name                string
	Last_Name                 string
	Email_Address             string
	Password                  string
	The_Country               Country
	Zip_or_Postal_Code        string
	City                      string
	State                     string
	Phone_Number              string
	Sobriety_Date             time.Time
	Birth_Date                time.Time
	Member_Of                 Fellowship
	YPAA_Committee            string
	Any_Special_Needs         SpecialNeeds
	Any_Service_Opportunities ServiceOpportunities
	Comments                  string
}

func PostRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprint(w, req.Form)
}

func GetRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	// signup_form.tmpl just needs a {{ .csrfField }} template tag for
	// csrf.TemplateField to inject the CSRF token into. Easy!
	t, _ := template.ParseFiles("signup_form.tmpl")
	t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(req),
		"countries":      Countries,
	})
}

type configuration struct {
	SiteName     string
	SiteDomain   string
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	ProjectID    string
	CSRF_Key     string
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
	router := CreateHandler(ContextHanderToHttpHandler)
	// TODO: comment out Dev and uncomment Live
	// Dev:
	csrfProtectedRouter := csrf.Protect([]byte(config.CSRF_Key), csrf.Secure(false))(router)
	// Live:
	//csrfProtectedRouter := csrf.Protect([]byte(config.CSRF_Key))(router)
	http.Handle("/", csrfProtectedRouter)
}
