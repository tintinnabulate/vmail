package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	c "github.com/smartystreets/goconvey/convey"

	"golang.org/x/net/context"
	//"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

var emailBody2 = `
Code: %s

Yours randomly,
Bert.
`

//func TestComposeVerificationEmail(t *testing.T) {
//	code := "abcde"
//	email := "foo@bar.baz"
//	want := &mail.Message{
//		Sender:  "[DONUT REPLY] Admin <donotreply@seraphic-lock-199316.appspotmail.com>",
//		To:      []string{email},
//		Subject: "Your verification code",
//		Body:    fmt.Sprintf(emailBody2, code),
//	}
//	if msg := ComposeVerificationEmail(email, code); !reflect.DeepEqual(msg, want) {
//		t.Errorf("composeMessage() = %+v, want %+v", msg, want)
//	}
//}

///

func CreateContextHandlerToHttpHandler(ctx context.Context) ContextHandlerToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("context: %v\n", ctx)
			fmt.Printf("w: %v\n", w)
			fmt.Printf("r: %v\n", r)
			fmt.Printf("f: %v\n", f)
			f(ctx, w, r)
		}
	}
}

//func TestVerifySignupEndpoint(t *testing.T) {
//	LoadConfig()
//
//	ctx, _, _ := aetest.NewContext()
//
//	c.Convey("When you want to do foo", t, func() {
//		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
//		record := httptest.NewRecorder()
//
//		req, err := http.NewRequest("GET", "/verify/38739873", nil)
//		c.So(err, c.ShouldBeNil)
//
//		c.Convey("It should return a 200 response", func() {
//			r.ServeHTTP(record, req)
//			c.So(record.Code, c.ShouldEqual, 200)
//			c.So(fmt.Sprint(record.Body), c.ShouldEqual, "foo hoi")
//		})
//	})
//}

//func TestCreateSignupEndpoint(t *testing.T) {
//	LoadConfig()
//
//	ctx, _, _ := aetest.NewContext()
//
//	c.Convey("When you want to do foo", t, func() {
//		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
//		record := httptest.NewRecorder()
//
//		req, err := http.NewRequest("POST", "/signup/foo@bar.com", nil)
//		c.So(err, c.ShouldBeNil)
//
//		c.Convey("It should return a 200 response", func() {
//			r.ServeHTTP(record, req)
//			c.So(record.Code, c.ShouldEqual, 200)
//			c.So(fmt.Sprint(record.Body), c.ShouldEqual, "foo hoi")
//		})
//	})
//}

func TestMonkeys(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/monkeys/dong", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			//c.So(fmt.Sprint(record.Body), c.ShouldEqual, "banana: dong")
		})
	})
}
