package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	c "github.com/smartystreets/goconvey/convey"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	//"google.golang.org/appengine"
)

func CreateContextHandlerToHttpHandler(ctx context.Context) ContextHandlerToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(ctx, w, r)
		}
	}
}

func TestPostRegistrationHandler(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When user tries to register with an unverified email address", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/register", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response, but suggest /signup", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, "Please sign up first /signup")
		})
	})
}

func TestGetRegistrationHandler(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When user visits the registration page", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/register", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, "hoi")
		})
	})
}
