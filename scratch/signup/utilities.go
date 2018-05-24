/*
	Implementation Note:
		None.
	Filename:
		utilities.go
*/

package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// CheckErr is a utility function for killing the app on the event of a non-nil error
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// RandToken generates a random token, for use in verification codes
func RandToken() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	CheckErr(err)
	return fmt.Sprintf("%x", b)
}

// HandlerFunc is our Standard http handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// ContextHandlerFunc is our context.Context http handler
type ContextHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// ContextHandlerToHandlerHOF is our Higher order function
// for changing a HandlerFunc to a ContextHandlerFunc, usually creating the context.Context along the way.
type ContextHandlerToHandlerHOF func(f ContextHandlerFunc) HandlerFunc

// ContextHandlerToHTTPHandler Creates a new Context and uses it when calling f
func ContextHandlerToHTTPHandler(f ContextHandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		f(ctx, w, r)
	}
}
