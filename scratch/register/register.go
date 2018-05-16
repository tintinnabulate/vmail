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

type Country int

const (
	United_States Country = iota + 1
	Canada
	Afghanistan
	Albania
	Algeria
	American_Samoa
	Andorra
	Angola
	Anguilla
	Antarctica
	Antigua_and_or_Barbuda
	Argentina
	Armenia
	Aruba
	Australia
	Austria
	Azerbaijan
	Bahamas
	Bahrain
	Bangladesh
	Barbados
	Belarus
	Belgium
	Belize
	Benin
	Bermuda
	Bhutan
	Bolivia
	Bosnia_and_Herzegovina
	Botswana
	Bouvet_Island
	Brazil
	British_lndian_Ocean_Territory
	Brunei_Darussalam
	Bulgaria
	Burkina_Faso
	Burundi
	Cambodia
	Cameroon
	Cape_Verde
	Cayman_Islands
	Central_African_Republic
	Chad
	Chile
	China
	Christmas_Island
	Cocos_Keeling_Islands
	Colombia
	Comoros
	Congo
	Cook_Islands
	Costa_Rica
	Croatia_Hrvatska
	Cuba
	Cyprus
	Czech_Republic
	Denmark
	Djibouti
	Dominica
	Dominican_Republic
	East_Timor
	Ecudaor
	Egypt
	El_Salvador
	Equatorial_Guinea
	Eritrea
	Estonia
	Ethiopia
	Falkland_Islands_Malvinas
	Faroe_Islands
	Fiji
	Finland
	France
	France_Metropolitan
	French_Guiana
	French_Polynesia
	French_Southern_Territories
	Gabon
	Gambia
	Georgia
	Germany
	Ghana
	Gibraltar
	Greece
	Greenland
	Grenada
	Guadeloupe
	Guam
	Guatemala
	Guinea
	Guinea_Bissau
	Guyana
	Haiti
	Heard_and_Mc_Donald_Islands
	Honduras
	Hong_Kong
	Hungary
	Iceland
	India
	Indonesia
	Iran_Islamic_Republic_of
	Iraq
	Ireland
	Israel
	Italy
	Ivory_Coast
	Jamaica
	Japan
	Jordan
	Kazakhstan
	Kenya
	Kiribati
	Korea_Democratic_Peoples_Republic_of
	Korea_Republic_of
	Kuwait
	Kyrgyzstan
	Lao_Peoples_Democratic_Republic
	Latvia
	Lebanon
	Lesotho
	Liberia
	Libyan_Arab_Jamahiriya
	Liechtenstein
	Lithuania
	Luxembourg
	Macau
	Macedonia
	Madagascar
	Malawi
	Malaysia
	Maldives
	Mali
	Malta
	Marshall_Islands
	Martinique
	Mauritania
	Mauritius
	Mayotte
	Mexico
	Micronesia_Federated_States_of
	Moldova_Republic_of
	Monaco
	Mongolia
	Montserrat
	Morocco
	Mozambique
	Myanmar
	Namibia
	Nauru
	Nepal
	Netherlands
	Netherlands_Antilles
	New_Caledonia
	New_Zealand
	Nicaragua
	Niger
	Nigeria
	Niue
	Norfork_Island
	Northern_Mariana_Islands
	Norway
	Oman
	Pakistan
	Palau
	Panama
	Papua_New_Guinea
	Paraguay
	Peru
	Philippines
	Pitcairn
	Poland
	Portugal
	Puerto_Rico
	Qatar
	Reunion
	Romania
	Russian_Federation
	Rwanda
	Saint_Kitts_and_Nevis
	Saint_Lucia
	Saint_Vincent_and_the_Grenadines
	Samoa
	San_Marino
	Sao_Tome_and_Principe
	Saudi_Arabia
	Senegal
	Seychelles
	Sierra_Leone
	Singapore
	Slovakia
	Slovenia
	Solomon_Islands
	Somalia
	South_Africa
	South_Georgia_South_Sandwich_Islands
	Spain
	Sri_Lanka
	St__Helena
	St__Pierre_and_Miquelon
	Sudan
	Suriname
	Svalbarn_and_Jan_Mayen_Islands
	Swaziland
	Sweden
	Switzerland
	Syrian_Arab_Republic
	Taiwan
	Tajikistan
	Tanzania_United_Republic_of
	Thailand
	Togo
	Tokelau
	Tonga
	Trinidad_and_Tobago
	Tunisia
	Turkey
	Turkmenistan
	Turks_and_Caicos_Islands
	Tuvalu
	Uganda
	Ukraine
	United_Arab_Emirates
	United_Kingdom
	United_States_minor_outlying_islands
	Uruguay
	Uzbekistan
	Vanuatu
	Vatican_City_State
	Venezuela
	Vietnam
	Virigan_Islands_British
	Virgin_Islands_U_S_
	Wallis_and_Futuna_Islands
	Western_Sahara
	Yemen
	Yugoslavia
	Zaire
	Zambia
	Zimbabwe
)

type Fellowship int

const (
	AA Fellowship = iota + 1
	AlAnon
	Alateen
	Other
)

type SpecialNeeds int

const (
	Deaf_or_Hard_of_Hearing SpecialNeeds = iota + 1
	Wheelchair_Access
	Translation_Services
)

type ServiceOpportunities int

const (
	Outreach ServiceOpportunities = iota + 1
	Service
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
