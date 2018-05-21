% Registration ICYfication
% Design Decisions - Issue 0.1, Draft
% \today{}
\newpage{}

# Design decisions

## Overview

I plan to work using the V-model from Software Engineering, developing using
Boehm's spiral.

My approach for user acceptance testing is to run through the User Stories (see
document entitled 'User Stories') and verify that I can carry-out each story
and get the expected result.

My approach for implementation (code) will be to use Test-First Development
(TFD), writing a test, running the code (showing that the test fails because
there is no code yet), and then writing the code to make the test pass.

## Implementation details

For the first release, I feel most comfortable writing the registration
application in Flask, with SQLAlchemy backend, running on Google App Engine. I
used this to make the Bath registration, and I can borrow heavily from this.

Once the first release is done, I can have people test it, and learn from it's
use. These lessons can inform the design of the second release.

I plan to write the second release mayyyyybe in Go, running on Google App Engine

Advantages of writing application in Go:

1. Go is fast (better than Python)
2. Strongly-typed, with lots of nice compile-time checks (better than Python)
3. Nice built-in testing framework
4. Extensive documentation about app development on Google App Engine exists

Disadvantages:

1. A lot fewer people understand Go than do Python. This is not to be
   underestimated - I want to be able to hand over the Registration for
   maintenance, with little-to-no input from myself.

## Stream of conscience

I want to just document some design ideas here... in an unstructured fashion...

Lets think through how the new convention host committee will use the registration site...
When they win the bid, I want to be able to just click a button and they get a registration page ready to take payments.

That form wants to be as simple as possible.

1. Admin enters the convention cost, and host city into a configuration file.
1. Admin pushes a 'Deploy' button.
1. Site is configured to ...
1. Registrant visits [https://eurypaa.org/signup]()
1. Registrant enters Email Address: xxx@xxx.com, clicks [SUBMIT]
1. Registrant gets verification link emailed to xxx@xxx.com 
1. Registrant clicks on verification link
1. Registrant is verified
1. Registrant is redirected to [https://eurypaa.org/register]()

* Someone must create an account in order to register
* A User object represents an account
* When a User registers, a Registration is created against that user, for that particular convention.
* Payment must be taken to create a Registration
* Convention is also an object


* We will need some facility for forgotten / resetting passwords... :(
* We will need some facility for the user modifying User details


Database Objects:

* User: 

```
type User struct {
    User_ID            int
    Creation_Date      time.Time
	First_Name         string
	Last_Name          string
	Email_Address      string
    Stripe_Customer_ID string
	Password           string
	Conf_Password      string
	The_Country        Country
	City               string
	Zip_or_Postal_Code string
	Sobriety_Date      time.Time
	Member_Of          []Fellowship
	Any_Special_Needs  []SpecialNeed
}
```

* Registration:

```
type Registration struct {
    Convention_ID       int
    Creation_Date       time.Time
    Stripe_Charge_ID    string
}
```

* Convention 

```
type Convention struct {
    Convention_ID   int
    Creation_Date   time.Time
    Year            int
    Country         EURYPAA_Country
    Cost            int
    Currency_Code   string
    Start_Date      time.Time
    End_Date        time.Time
    Venue           string
    Hotel           string
    Venue_Is_Hotel  bool
}
```

* Tshirt? `{size TshirtSize, ...}`

# Change history

Date | Issue | Note
---|---|---
2017-03-09 | 0.1, Draft | Initial draft |
