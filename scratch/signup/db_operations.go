/*
	Implementation Note:
		None.

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
	SiteCode          string    `datastore:"site_code"`
	IsVerified        bool      `datastore:"verified"`
	id                int64     // The integer ID used in the datastore.
}

type Site struct {
	CreationTimestamp time.Time `datastore:"created"`
	SiteName          string    `datastore:"site_name"`
	Code              string    `datastore:"code"`
	RootURL           string    `datastore:"root_url"`
	VerifiedURL       string    `datastore:"verified_url"`
	id                int64     // The integer ID used in the datastore.
}

// AddSignup adds a signup with the given verification code to the datastore,
// returning the key of the newly created entity.
func AddSignup(ctx context.Context, siteCode, email, code string) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)
	signup := &Signup{
		CreationTimestamp: time.Now(),
		Email:             email,
		VerificationCode:  code,
		SiteCode:          siteCode,
		IsVerified:        false,
	}
	k, err := datastore.Put(ctx, key, signup)
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
	signup := new(Signup)
	if err := datastore.Get(ctx, key, signup); err != nil {
		return true, nil
	}
	return false, nil
}

// MarkVerified marks the signup as verified with the given ID.
func MarkVerified(ctx context.Context, code string) error {
	// Create a key using the given integer ID.
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)
	signup := new(Signup)
	// In a transaction load each signup, set verified to true and store.
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		if err := datastore.Get(tx, key, signup); err != nil {
			return errors.New("no such verification code")
		}
		if signup.IsVerified {
			return errors.New("signup already verified")
		}
		signup.IsVerified = true
		_, err := datastore.Put(tx, key, signup)
		return err
	}, nil)
	return err
}

// GetSite gets a site matching siteCode
func GetSite(ctx context.Context, siteCode string) (Site, error) {
	key := datastore.NewKey(ctx, "Site", siteCode, 0, nil)
	var site Site
	if err := datastore.Get(ctx, key, &site); err != nil {
		return site, fmt.Errorf("GetSite: datastore.Get: %v", err)
	}
	return site, nil
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

// AddSite adds a site with the given verification code to the datastore,
// returning the key of the newly created entity.
// Should only be called during testing.
func AddSite(ctx context.Context, siteName, siteCode, rootURL string) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "Site", siteCode, 0, nil)
	site := &Site{
		CreationTimestamp: time.Now(),
		SiteName:          siteName,
		Code:              siteCode,
		RootURL:           rootURL,
		VerifiedURL:       rootURL,
	}
	k, err := datastore.Put(ctx, key, site)
	return k, err
}
