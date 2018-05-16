package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

// Creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFuncs.
func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup", f(GetSignupHandler)).Methods("GET")
	appRouter.HandleFunc("/signup", f(PostSignupHandler)).Methods("POST")
	appRouter.HandleFunc("/register", f(GetRegistrationHandler)).Methods("GET")
	appRouter.HandleFunc("/register", f(PostRegistrationHandler)).Methods("POST")

	return appRouter
}

func GetSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	// TODO: we can load templates nicer than this - do it once, and globally
	t, err := template.New("signup_form.tmpl").ParseFiles("signup_form.tmpl")
	CheckErr(err)
	t.ExecuteTemplate(w,
		"signup_form.tmpl",
		map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func PostSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	CheckErr(err)
	var signup Signup
	err = schemaDecoder.Decode(&signup, req.PostForm)
	client := urlfetch.Client(ctx)
	resp, err := client.Post(fmt.Sprintf("%s/%s", config.SignupURL, signup.Email_Address), "", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)
}

func GetRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	t, err := template.New("register_form.tmpl").Funcs(funcMap).ParseFiles("register_form.tmpl")
	CheckErr(err)
	t.ExecuteTemplate(w,
		"register_form.tmpl",
		map[string]interface{}{
			"Countries":            Countries,
			"Fellowships":          Fellowships,
			"SpecialNeeds":         SpecialNeeds,
			"ServiceOpportunities": ServiceOpportunities,
			csrf.TemplateTag:       csrf.TemplateField(req),
		})
}

func PostRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	CheckErr(err)
	var registration Registration
	err = schemaDecoder.Decode(&registration, req.PostForm)
	CheckErr(err)
	// registration now holds our user
	// TODO:
	// 1. `resp, err := http.Get(fmt.Sprinf("signup_verifier.com/signup/%s", registration.Email_Address))`
	// 2. `if resp.Body != "{'success':true}" { redirect("/signup") }
	fmt.Fprint(w, registration)
}

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

type Signup struct {
	Email_Address string
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
	SignupURL    string
}

var (
	config        configuration
	schemaDecoder = schema.NewDecoder()
	funcMap       = template.FuncMap{"inc": func(i int) int { return i + 1 }}
)

func Config_Init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	CheckErr(err)
}

func SchemaDecoder_Init() {
	schemaDecoder.RegisterConverter(time.Time{}, TimeConverter)
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
