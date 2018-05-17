/*
	Implementation Note:
		All calls to `datastore.Put` should be followed by a `datastore.Get`.
		This forces the `Put` to store immediately when run locally, which is
		necessary for testing with `goapp test`.
		See more info here: https://stackoverflow.com/a/25075074

	Filename:
		db_operations.go
*/

package main

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// Signup is used to hold signup entries going into/coming out of the datastore
type Signup struct {
	CreationTimestamp time.Time `datastore:"created"`
	Email             string    `datastore:"email"`
	VerificationCode  string    `datastore:"code"`
	IsVerified        bool      `datastore:"verified"`
	id                int64     // The integer ID used in the datastore.
}

// AddSignup adds a signup with the given verification code to the datastore,
// returning the key of the newly created entity.
func AddSignup(ctx context.Context, email, code string) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)
	signup := &Signup{
		CreationTimestamp: time.Now(),
		Email:             email,
		VerificationCode:  code,
		IsVerified:        false,
	}
	k, err := datastore.Put(ctx, key, signup)
	err2 := datastore.Get(ctx, k, &signup)
	CheckErr(err2)
	return k, err
}

// IsSignupVerified checks the database to see if an email address exits and is verified.
func IsSignupVerified(ctx context.Context, email string) (bool, error) {
	q := datastore.NewQuery("Signup").
		Filter("email =", email).
		Filter("verified =", true)
	var signups []Signup
	if _, err := q.GetAll(ctx, &signups); err != nil {
		return false, err
	}
	if len(signups) < 1 {
		return false, nil
	}
	return true, nil
}

// IsCodeAvailable checks the database to see if code is free to use.
func IsCodeAvailable(ctx context.Context, code string) (bool, error) {
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)
	var signup Signup
	if err := datastore.Get(ctx, key, &signup); err != nil {
		return true, nil
	}
	return false, nil
}

// MarkVerified marks the signup as verified with the given ID.
func MarkVerified(ctx context.Context, code string) error {
	// Create a key using the given integer ID.
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)
	var signup Signup
	// In a transaction load each signup, set verified to true and store.
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		if err := datastore.Get(tx, key, &signup); err != nil {
			return errors.New("no such verification code")
		}
		if signup.IsVerified {
			return errors.New("signup already verified")
		}
		signup.IsVerified = true
		_, err := datastore.Put(tx, key, &signup)
		return err
	}, nil)
	err2 := datastore.Get(ctx, key, &signup)
	CheckErr(err2)
	return err
}

// GetSignupCode gets the signup code matching the given email address.
// This should only be called during testing.
func GetSignupCode(ctx context.Context, email string) (string, error) {
	q := datastore.NewQuery("Signup").Filter("email =", email)
	var signups []Signup
	if _, err := q.GetAll(ctx, &signups); err != nil {
		return "", err
	}
	if len(signups) < 1 {
		return "", fmt.Errorf("Email not in database: %s", email)
	}
	return signups[0].VerificationCode, nil
}
