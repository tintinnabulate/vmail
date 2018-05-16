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

// FIXME: this is a horrendous hack to make up for enums starting at 1 but loops starting at 0
func amendPostForm(req *http.Request) {
	x, _ := strconv.Atoi(req.PostForm["The_Country"][0])
	req.PostForm["The_Country"][0] = fmt.Sprint(x + 1)
	for i, el := range req.PostForm["Any_Special_Needs"] {
		x, _ = strconv.Atoi(el)
		req.PostForm["Any_Special_Needs"][i] = fmt.Sprint(x + 1)
	}
	for i, el := range req.PostForm["Any_Service_Opportunities"] {
		x, _ = strconv.Atoi(el)
		req.PostForm["Any_Service_Opportunities"][i] = fmt.Sprint(x + 1)
	}
	for i, el := range req.PostForm["Member_Of"] {
		x, _ = strconv.Atoi(el)
		req.PostForm["Member_Of"][i] = fmt.Sprint(x + 1)
	}
}

func PostRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	checkErr(err)
	var registration Registration

	// FIXME
	amendPostForm(req)

	err = schemaDecoder.Decode(&registration, req.PostForm)
	//checkErr(err) // TODO: schema can't handle gorilla CSRT token... how to handle?
	fmt.Fprint(w, registration)
}

func GetRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("signup_form.tmpl")
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
}

var (
	config        configuration
	appRouter     mux.Router
	schemaDecoder = schema.NewDecoder()
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func timeConverter(value string) reflect.Value {
	tstamp, _ := strconv.ParseInt(value, 10, 64)
	return reflect.ValueOf(time.Unix(tstamp, 0))
}

func LoadConfig() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	checkErr(err)
	schemaDecoder.RegisterConverter(time.Time{}, timeConverter)
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
