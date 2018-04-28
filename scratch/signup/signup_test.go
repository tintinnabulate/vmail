package main

import (
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/appengine/mail"
)

var emailBody2 = `
Code: %s

Yours randomly,
Bert.
`

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
