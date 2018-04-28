package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"google.golang.org/appengine/mail"
)

var emailBody2 = `
Code: %s

Yours randomly,
Bert.
`

func TestMain(m *testing.M) {
	Initialise()

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	appRouter.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestComposeVerificationEmail(t *testing.T) {
	code := "abcde"
	email := "foo@bar.baz"
	want := &mail.Message{
		Sender:  "[DONUT REPLY] Admin <donotreply@seraphic-lock-199316.appspotmail.com>",
		To:      []string{email},
		Subject: "Your verification code",
		Body:    fmt.Sprintf(emailBody2, code),
	}
	if msg := ComposeVerificationEmail(email, code); !reflect.DeepEqual(msg, want) {
		t.Errorf("composeMessage() = %+v, want %+v", msg, want)
	}
}

//func TestCreateSignupEndpoint(t *testing.T) {
//	req, _ := http.NewRequest("POST", "/signup/justwanttouseappspot@gmail.com", nil)
//	response := executeRequest(req)
//	checkResponseCode(t, http.StatusOK, response.Code)
//	if body := response.Body.String(); body != "[]" {
//		t.Errorf("Expected an empty array. Got %s", body)
//	}
//}

func TestVerifySignup(t *testing.T) {
	req, _ := http.NewRequest("GET", "/verify/4839202", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}
