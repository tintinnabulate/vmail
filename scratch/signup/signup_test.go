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

func CreateContextHandlerToHttpHandler(ctx context.Context) ContextHandlerToHandlerHOF {
	return func(f ContextHandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f(ctx, w, r)
		}
	}
}

func TestCreateSignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
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

func TestCreateAndVerifyAndCheckSignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()
		record2 := httptest.NewRecorder()
		record3 := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/signup/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {

			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
`)

			code, _ := GetSignupCode(ctx, "lolz")

			req2, err2 := http.NewRequest("GET", fmt.Sprintf("/verify/%s", code), nil)
			c.So(err2, c.ShouldBeNil)

			c.Convey("Should be 200", func() {

				r.ServeHTTP(record2, req2)
				c.So(record2.Code, c.ShouldEqual, 200)
				c.So(fmt.Sprint(record2.Body), c.ShouldEqual, fmt.Sprintf(`{"code":"%s","Success":true,"Note":""}
`, code))

				req3, err3 := http.NewRequest("GET", "/signup/lolz", nil)
				c.So(err3, c.ShouldBeNil)

				c.Convey("Should be 200", func() {
					r.ServeHTTP(record3, req3)
					c.So(record3.Code, c.ShouldEqual, 200)
					c.So(fmt.Sprint(record3.Body), c.ShouldEqual, `{"address":"lolz","success":false,"note":""}
`)

				})
			})
		})
	})
}

func TestVerifySignupEndpoint(t *testing.T) {
	LoadConfig()

	ctx, _, _ := aetest.NewContext()

	c.Convey("When you want to do foo", t, func() {
		r := CreateHandler(CreateContextHandlerToHttpHandler(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/verify/lolz", nil)
		c.So(err, c.ShouldBeNil)

		c.Convey("It should return a 200 response", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, 200)
			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"code":"lolz","Success":false,"Note":"no such verification code"}
`)
		})
	})
}

// TODO test adding a valid signup and look for that code
