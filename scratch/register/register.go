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
	"github.com/stripe/stripe-go"
	stripeClient "github.com/stripe/stripe-go/client"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

func CreateHandler(f ContextHandlerToHandlerHOF) *mux.Router {
	appRouter := mux.NewRouter()
	appRouter.HandleFunc("/signup", f(GetSignupHandler)).Methods("GET")
	appRouter.HandleFunc("/signup", f(PostSignupHandler)).Methods("POST")
	appRouter.HandleFunc("/register", f(GetRegistrationHandler)).Methods("GET")
	appRouter.HandleFunc("/register", f(PostRegistrationHandler)).Methods("POST")
	appRouter.HandleFunc("/charge", f(GetRegistrationPaymentHandler)).Methods("GET")
	appRouter.HandleFunc("/charge", f(PostRegistrationPaymentHandler)).Methods("POST")

	return appRouter
}

func GetSignupHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w,
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
	t, err := template.New("registration_form.tmpl").Funcs(funcMap).ParseFiles("registration_form.tmpl")
	CheckErr(err)
	t.ExecuteTemplate(w,
		"registration_form.tmpl",
		map[string]interface{}{
			"Countries":            Countries,
			"Fellowships":          Fellowships,
			"SpecialNeeds":         SpecialNeeds,
			"ServiceOpportunities": ServiceOpportunities,
			csrf.TemplateTag:       csrf.TemplateField(req),
		})
}

func PostRegistrationHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var registration Registration
	var signup Signup
	err := req.ParseForm()
	CheckErr(err)
	err = schemaDecoder.Decode(&registration, req.PostForm)
	CheckErr(err)
	client := urlfetch.Client(ctx)
	resp, err := client.Get(fmt.Sprintf("%s/%s", config.SignupURL, registration.Email_Address))
	CheckErr(err)
	json.NewDecoder(resp.Body).Decode(&signup)
	if signup.Success {
		fmt.Fprint(w, "You may proceed %v", registration)
	} else {
		fmt.Fprint(w, "I'm sorry, you need to sign up first. Go to /signup")
	}
}

func GetRegistrationPaymentHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	tmpl := templates.Lookup("stripe.tmpl")
	tmpl.Execute(w,
		map[string]interface{}{
			"Key":            publishableKey,
			csrf.TemplateTag: csrf.TemplateField(req),
		})
}

func PostRegistrationPaymentHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	customerParams := &stripe.CustomerParams{Email: r.Form.Get("stripeEmail")}
	customerParams.SetSource(r.Form.Get("stripeToken"))

	httpClient := urlfetch.Client(ctx)
	sc := stripeClient.New(stripe.Key, stripe.NewBackends(httpClient))

	newCustomer, err := sc.Customers.New(customerParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   500,
		Currency: "usd",
		Desc:     "Sample Charge",
		Customer: newCustomer.ID,
	}
	charge, err := sc.Charges.New(chargeParams)
	if err != nil {
		fmt.Fprintf(w, "Could not process payment: %v", err)
	}
	fmt.Fprintf(w, "Completed payment: %v", charge.ID)
}

type Registration struct {
	First_Name         string
	Last_Name          string
	Email_Address      string
	Password           string
	Conf_Password      string
	The_Country        Country
	City               string
	Zip_or_Postal_Code string
	Sobriety_Date      time.Time
	Member_Of          []Fellowship
	Any_Special_Needs  []SpecialNeed
}

type Signup struct {
	Email_Address string `json:"address"`
	Success       bool   `json:"success"`
	Note          string `json:"note"`
}

type configuration struct {
	SiteName             string
	SiteDomain           string
	SMTPServer           string
	SMTPUsername         string
	SMTPPassword         string
	ProjectID            string
	CSRF_Key             string
	IsLiveSite           bool
	SignupURL            string
	StripePublishableKey string
	StripeSecretKey      string
}

var (
	config         configuration
	schemaDecoder  = schema.NewDecoder()
	funcMap        = template.FuncMap{"inc": func(i int) int { return i + 1 }}
	publishableKey string
	templates      = template.Must(template.ParseGlob("views/*.tmpl"))
)

func ConfigInit() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	config = configuration{}
	err := decoder.Decode(&config)
	CheckErr(err)
}

func SchemaDecoderInit() {
	schemaDecoder.RegisterConverter(time.Time{}, TimeConverter)
	schemaDecoder.IgnoreUnknownKeys(true)
}

func RouterInit() {
	// TODO: https://youtu.be/xyDkyFjzFVc?t=1308
	router := CreateHandler(ContextHandlerToHttpHandler)
	csrfProtector := csrf.Protect(
		[]byte(config.CSRF_Key),
		csrf.Secure(config.IsLiveSite))
	csrfProtectedRouter := csrfProtector(router)
	http.Handle("/", csrfProtectedRouter)
}

func StripeInit() {
	publishableKey = config.StripePublishableKey
	stripe.Key = config.StripeSecretKey
}

func init() {
	ConfigInit()
	SchemaDecoderInit()
	RouterInit()
	StripeInit()
}
