package main

import (
	"time"
)

type User struct {
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
	Member_Of                 []Fellowship
	YPAA_Committee            string
	Any_Special_Needs         []SpecialNeed
	Any_Service_Opportunities []ServiceOpportunity
}
