package main

import (
	"errors"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

type Signup struct {
	CreationTimestamp time.Time `datastore:"created"`
	Email             string    `datastore:"email"`
	VerificationCode  string    `datastore:"code"`
	IsVerified        bool      `datastore:"verified"`
	id                int64     // The integer ID used in the datastore.
}

// AddSignup adds a signup with the given verification code to the datastore,
// returning the key of the newly created entity.
func AddSignup(r *http.Request, email, code string) (*datastore.Key, error) {
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)
	signup := &Signup{
		CreationTimestamp: time.Now(),
		Email:             email,
		VerificationCode:  code,
		IsVerified:        false,
	}
	return datastore.Put(ctx, key, signup)
}

// MarkDone marks the signup as verified with the given ID.
func MarkVerified(r *http.Request, code string) error {
	ctx := appengine.NewContext(r)
	// Create a key using the given integer ID.
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)

	// In a transaction load each signup, set verified to true and store.
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		var signup Signup
		if err := datastore.Get(tx, key, &signup); err != nil {
			return errors.New("no such verification code")
		}
		if signup.IsVerified {
			return errors.New("signup already verified")
		} else {
			signup.IsVerified = true
			_, err := datastore.Put(tx, key, &signup)
			return err
		}
	}, nil)
	return err
}
