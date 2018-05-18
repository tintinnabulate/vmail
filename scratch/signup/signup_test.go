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

// TestCreateSignupEndpoint tests that we can create a signup
func TestCreateSignupEndpoint(t *testing.T) {
	LoadConfig()

	inst, _ := aetest.NewInstance(
		&aetest.Options{
			StronglyConsistentDatastore: true,
		})
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
	}
	ctx := appengine.NewContext(req)

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/signup/lolz", nil)
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

	inst, _ := aetest.NewInstance(
		&aetest.Options{
			StronglyConsistentDatastore: true,
		})
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
	}
	ctx := appengine.NewContext(req)

	c.Convey("When creating a signup for email address 'lolz'", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()  // records the 'create signup' response
		record2 := httptest.NewRecorder() // records the 'verify signup' response
		record3 := httptest.NewRecorder() // records the 'check signup is verified' repsonse

		req, err := http.NewRequest("POST", "/signup/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should succeed", func() {

			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)

			// Look up code sent to 'lolz'
			code, _ := GetSignupCode(ctx, "lolz")

			req2, err2 := http.NewRequest("GET", fmt.Sprintf("/verify/%s", code), nil)
			c.So(err2, c.ShouldBeNil)

			c.Convey("Verifying the code sent to 'lolz' should succeed", func() {

				r.ServeHTTP(record2, req2)
				c.So(record2.Code, c.ShouldEqual, 200)
				c.So(fmt.Sprint(record2.Body), c.ShouldEqual, fmt.Sprintf(`{"code":"%s","success":true,"note":""}
`, code))

				req3, err3 := http.NewRequest("GET", "/signup/lolz", nil)
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

	inst, _ := aetest.NewInstance(
		&aetest.Options{
			StronglyConsistentDatastore: true,
		})
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
	}
	ctx := appengine.NewContext(req)

	c.Convey("When you try and verify a non-existent code", t, func() {
		r := CreateHandler(CreateContextHandlerToHTTPHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/verify/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response, but fail", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"code":"lolz","success":false,"note":"no such verification code"}
`)
		})
	})
}
