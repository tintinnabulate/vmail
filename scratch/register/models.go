package main

import (
	"time"
)

type UserModel struct {
	User_ID            int // pk
	Creation_Date      time.Time
	First_Name         string
	Last_Name          string
	Email_Address      string
	Password           string
	Conf_Password      string
	The_CountryType    CountryType
	City               string
	Sobriety_Date      time.Time
	Member_Of          []Fellowship
	Stripe_Customer_ID string
}

type RegistrationModel struct {
	Registration_ID  int // pk
	User_ID          int // fk
	Convention_ID    int // fk
	Creation_Date    time.Time
	Stripe_Charge_ID string
}

type ConventionModel struct {
	Convention_ID     int // pk
	Creation_Date     time.Time
	Year              int
	The_CountryType   CountryType
	Cost              int
	Currency_Code     string
	Start_Date        time.Time
	End_Date          time.Time
	Hotel             string
	Hotel_Is_Venue    bool
	Venue             string
	Stripe_Product_ID string
}
