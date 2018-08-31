package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	c "github.com/smartystreets/goconvey/convey"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
)

func CreateContextHandlerToHTTPHandler(ctx context.Context) ContextHandlerToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(ctx, w, r)
		}
	}
}

func getContext() (context.Context, aetest.Instance) {
	inst, _ := aetest.NewInstance(
		&aetest.Options{
			StronglyConsistentDatastore: true,
			// SuppressDevAppServerLog:     true,
		})
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
	}
	ctx := appengine.NewContext(req)
	return ctx, inst
}

// TestCreateSignupEndpoint tests that we can create a signup
func TestCreateSignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, inst := getContext()
	defer inst.Close()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/signup/foo/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)
		})
	})
}

// TestCreateAndVerifyAndCheckSignupEndpoint tests that we can create a signup, verify it, and then check that it is verified
func TestCreateAndVerifyAndCheckSignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, inst := getContext()
	defer inst.Close()

	c.Convey("When creating a signup for email address 'lolz'", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()  // records the 'create signup' response
		record2 := httptest.NewRecorder() // records the 'verify signup' response
		record3 := httptest.NewRecorder() // records the 'check signup is verified' repsonse

		req, err := http.NewRequest("POST", "/signup/foo/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should succeed", func() {

			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)

			// Look up code sent to 'lolz'
			code, _ := GetSignupCode(ctx, "lolz")

			_, errk := AddSite(ctx, "FOOWEBSITE", "foo", "http://barnacles.com")
			CheckErr(errk)

			req2, err2 := http.NewRequest("GET", fmt.Sprintf("/verify/foo/%s", code), nil)
			c.So(err2, c.ShouldBeNil)

			c.Convey("Verifying the code sent to 'lolz' should succeed", func() {

				r.ServeHTTP(record2, req2)
				c.So(record2.Code, c.ShouldEqual, 302)
				c.So(fmt.Sprint(record2.Body), c.ShouldEqual, fmt.Sprint(`<a href="http://barnacles.com">Found</a>.

`))

				req3, err3 := http.NewRequest("GET", "/signup/foo/lolz", nil)
				c.So(err3, c.ShouldBeNil)

				c.Convey("Checking email 'lolz' is verified should succeed", func() {
					r.ServeHTTP(record3, req3)
					c.So(record3.Code, c.ShouldEqual, 200)
					c.So(fmt.Sprint(record3.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)

				})
			})
		})
	})
}

// TestVerifySignupEndpoint tests that verifying a non-existent code produces a JSON response where "success": false .
func TestVerifySignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, inst := getContext()
	defer inst.Close()

	c.Convey("When you try and verify a non-existent code", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/verify/foo/lolz", nil)
		c.So(err, c.ShouldBeNil)

		_, errk := AddSite(ctx, "FOOWEBSITE", "foo", "http://barnacles.com")
		CheckErr(errk)

		c.Convey("It should return a 200 response, but fail", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 303)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `<a href="http://barnacles.com">See Other</a>.

`)
		})
	})
}
