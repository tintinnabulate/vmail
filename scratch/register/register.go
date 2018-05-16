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
	CheckErr(err)
	var registration Registration
	err = schemaDecoder.Decode(&registration, req.PostForm)
	CheckErr(err)
	fmt.Fprint(w, registration)
}

func GetRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	t, err := template.New("signup_form.tmpl").Funcs(funcMap).ParseFiles("signup_form.tmpl")
	CheckErr(err)
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

// Creates my mux.Router. Uses f to convert ContextHandlerFunc's to HandlerFuncs.
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
