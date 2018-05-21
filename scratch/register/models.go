package main

import (
	"time"
)

type User struct {
	User_ID            int64 // pk
	Creation_Date      time.Time
	First_Name         string
	Last_Name          string
	Email_Address      string
	Password           string
	Conf_Password      string
	Country            CountryType
	City               string
	Sobriety_Date      time.Time
	Member_Of          []Fellowship
	Stripe_Customer_ID string
}

type Registration struct {
	Registration_ID  int64 // pk
	User_ID          int64 // fk
	Convention_ID    int64 // fk
	Creation_Date    time.Time
	Stripe_Charge_ID string
}

type Convention struct {
	Convention_ID     int64 // pk
	Creation_Date     time.Time
	Year              int
	Country           CountryType
	Cost              int
	Currency_Code     string
	Start_Date        time.Time
	End_Date          time.Time
	Hotel             string
	Hotel_Is_Venue    bool
	Venue             string
	Stripe_Product_ID string
}

type RegistrationForm struct {
	id            int64 // pk
	First_Name    string
	Last_Name     string
	Email_Address string
	Password      string
	Conf_Password string
	Country       CountryType
	City          string
	Sobriety_Date time.Time
	Member_Of     []Fellowship
}
