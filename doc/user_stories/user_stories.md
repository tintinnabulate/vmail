% Registration ICYfication
% User Stories - Issue 0.2, Draft
% \today{}
\newpage{}

# Stories

## [US-1-1] When user signs up

**Pre-requisites**

1. None

**Story**

1. Access the address "9.ep.org" from their web browser
2. Click 'Sign in / Sign up' (9.ep.org/sign_in GET)
3. Fill in form to sign up & submit (9.ep.org/sign_in POST)
4. Click validation link in email sent to user (9.ep.org/email_validation GET)
5. End.

## [US-2-1] When user registers

**Pre-requisites**

1. User has not registered

**Story**

1. Access the address "9.ep.org" from their web browser
2. Click 'Register'. If they are signed in, GOTO 3, else GOTO 4.
3. Fill in form & submit (9.ep.org/register/1 GET/POST)
4. Fill in form to 'Sign in / Sign up', and redirect to /register afterwards (9.ep.org/sign_in GET/POST)
5. Fill in form to pay (9.ep.org/register/2 GET/POST) - Stripe
6. Send user confirmation email
7. End.

## [US-3-1] When registered user tries to register again

**Pre-requisites**

1. User has an account
2. User is signed in
3. User has registered
4. End.

**Story**

1. User clicks 'Register'
2. Site lists that they have already registered, suggesting the 'My Account' options from US-10-1.
3. End.

## [US-4-1] When a non-signed up user tries to sign in

**Pre-requisites**

1. User has not signed up

**Story**

1. Access the address "9.ep.org" from their web browser
2. Click 'Sign in / Sign up' (9.ep.org/sign_in GET)
3. Fill in form to sign in & submit (9.ep.org/sign_in POST)
4. Display 'email address/password not recognised' page. (9.ep.org/sign_in GET)
5. End.

## [US-5-1] When registered user signs in

**Pre-requisites**

1. User is has signed up
2. User has registered

**Story**

1. Access the address "9.ep.org" from their web browser
2. Click 'Sign in / Sign up' (9.ep.org/sign_in GET)
3. Fill in form to sign in & submit (9.ep.org/sign_in POST)
4. Take them back to the page they were on
5. End.

## [US-6-1] When user wants a refund (?)

Do we want to support this?

## [US-7-1] When user clicks 'Transfer registration to another person'

i.e. "I can't make it any more, but my friend can!"

## [US-8-1] When user clicks 'Buy Scholarship for a specific person'

e.g. "I'm paying for a sponsee"

## [US-9-1] When user clicks 'Buy Scholarship for anyone in need'

e.g. "I'm donating 1 registration for someone in need"

## [US-10-1] When user clicks 'My Account'

**Pre-requisites**

1. User is signed in
2. User has registered

**Story**

1. Show 'My Account' page options:
	1. Transfer Registration
	2. Donate Scholarship for a specific person
	3. Donate Scholarship for anyone in need
2. End.

## [US-11-1] When a person clicks 'Delete account'

1. Display if they have any registrations still paid for, for conventions not-yet-elapsed, and display a warning/suggestion to transfer registration
2. Display message saying their details will be removed from storage
3. Delete the account, removing all details from DB.
4. End.

# Change history

Date | Issue | Note
---|---|---
2017-03-09 | 0.1, Draft | Initial draft |
2017-03-09 | 0.2, Draft | Formatting |
