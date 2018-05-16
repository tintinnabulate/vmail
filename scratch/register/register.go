package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// Maybe you want to use github.com/gorilla/schema
type Registration struct {
	First_Name                string
	Last_Name                 string
	Email_Address             string
	Password                  string
	Conf_Password             string
	The_Country               Country
	Zip_or_Postal_Code        string
	City                      string
	State                     string
	Phone_Number              string
	Sobriety_Date             time.Time
	Birth_Date                time.Time
	Member_Of                 []Fellowship
	YPAA_Committee            string
	Any_Special_Needs         []SpecialNeed
	Any_Service_Opportunities []ServiceOpportunity
	Comments                  string
}

func PostRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	checkErr(err)
	var registration Registration
	err = schemaDecoder.Decode(&registration, req.PostForm)
	checkErr(err)
	fmt.Fprint(w, registration)
}

func GetRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	t, err := template.New("signup_form.tmpl").Funcs(funcMap).ParseFiles("signup_form.tmpl")
	checkErr(err)
	t.ExecuteTemplate(w,
		"signup_form.tmpl",
		map[string]interface{}{
			"Countries":            Countries,
			"Fellowships":          Fellowships,
			"SpecialNeeds":         SpecialNeeds,
			"ServiceOpportunities": ServiceOpportunities,
			csrf.TemplateTag:       csrf.TemplateField(req),
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
	IsLiveSite   bool
}

var (
	config        configuration
	appRouter     mux.Router
	schemaDecoder = schema.NewDecoder()
	funcMap       = template.FuncMap{"inc": func(i int) int { return i + 1 }}
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: this will need adapting to whatever format we request for Sobriety_Date and Birth_Date
func timeConverter(value string) reflect.Value {
	tstamp, _ := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(time.Unix(tstamp, 0))
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

func Config_Init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
}

func SchemaDecoder_Init() {
	schemaDecoder.RegisterConverter(time.Time{}, timeConverter)
	schemaDecoder.IgnoreUnknownKeys(true)
}

func Router_Init() {
	router := CreateHandler(ContextHanderToHttpHandler)
	csrfProtector := csrf.Protect(
		[]byte(config.CSRF_Key),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

func init() {
	Config_Init()
	SchemaDecoder_Init()
	Router_Init()
}
