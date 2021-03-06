package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	os.Exit(retCode)
}

func testSetup() {
	configInit("config.example.json")
	environmentInit()
}

// TestCreateSignupEndpoint tests that we can create a signup
func TestCreateSignupEndpoint(t *testing.T) {
	req, err := http.NewRequest("POST", "/signup/foo/bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateSignupEndpoint)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}
	expected := (`{"address":"bar","success":true,"note":""}
`)
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf(
			"unexpected body: got (%v) want (%v)",
			rr.Body.String(),
			expected,
		)
	}
}

//// TestCreateAndVerifyAndCheckSignupEndpoint tests that we can create a signup, verify it, and then check that it is verified
//func TestCreateAndVerifyAndCheckSignupEndpoint(t *testing.T) {
//	ctx := context.Background()
//	c.Convey("When creating a signup for email address 'lolz'", t, func() {
//		r := CreateHandler()
//		record := httptest.NewRecorder()  // records the 'create signup' response
//		record2 := httptest.NewRecorder() // records the 'verify signup' response
//		record3 := httptest.NewRecorder() // records the 'check signup is verified' repsonse
//
//		req, err := http.NewRequest("POST", "/signup/foo/lolz", nil)
//		c.So(err, c.ShouldBeNil)
//
//		c.Convey("It should succeed", func() {
//
//			r.ServeHTTP(record, req)
//			c.So(record.Code, c.ShouldEqual, 200)
//			c.So(fmt.Sprint(record.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
//`)
//
//			// Look up code sent to 'lolz'
//			code, _ := GetSignupCode(ctx, "lolz")
//
//			_, errk := AddSite(ctx, "FOOWEBSITE", "foo", "http://barnacles.com")
//			CheckErr(errk)
//
//			req2, err2 := http.NewRequest("GET", fmt.Sprintf("/verify/foo/%s", code), nil)
//			c.So(err2, c.ShouldBeNil)
//
//			c.Convey("Verifying the code sent to 'lolz' should succeed", func() {
//
//				r.ServeHTTP(record2, req2)
//				c.So(record2.Code, c.ShouldEqual, 302)
//				c.So(fmt.Sprint(record2.Body), c.ShouldEqual, fmt.Sprint(`<a href="http://barnacles.com">Found</a>.
//
//`))
//
//				req3, err3 := http.NewRequest("GET", "/signup/foo/lolz", nil)
//				c.So(err3, c.ShouldBeNil)
//
//				c.Convey("Checking email 'lolz' is verified should succeed", func() {
//					r.ServeHTTP(record3, req3)
//					c.So(record3.Code, c.ShouldEqual, 200)
//					c.So(fmt.Sprint(record3.Body), c.ShouldEqual, `{"address":"lolz","success":true,"note":""}
//`)
//
//				})
//			})
//		})
//	})
//}
//
//// TestVerifySignupEndpoint tests that verifying a non-existent code produces a JSON response where "success": false .
//func TestVerifySignupEndpoint(t *testing.T) {
//	ctx := context.Background()
//	c.Convey("When you try and verify a non-existent code", t, func() {
//		r := CreateHandler()
//		record := httptest.NewRecorder()
//
//		req, err := http.NewRequest("GET", "/verify/foo/lolz", nil)
//		c.So(err, c.ShouldBeNil)
//
//		_, errk := AddSite(ctx, "FOOWEBSITE", "foo", "http://barnacles.com")
//		CheckErr(errk)
//
//		c.Convey("It should return StatusNotFound, and fail", func() {
//			r.ServeHTTP(record, req)
//			c.So(record.Code, c.ShouldEqual, http.StatusNotFound)
//		})
//	})
//}
