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

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
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

// Site : stores each site we verify emails for
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
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("could not create datastore client: %v", err)
	}
	key := datastore.IncompleteKey("Signup", nil)
	signup := &Signup{
		CreationTimestamp: time.Now(),
		Email:             email,
		VerificationCode:  code,
		SiteCode:          siteCode,
		IsVerified:        false,
	}
	if _, err := client.Put(ctx, key, signup); err != nil {
		return nil, fmt.Errorf("could not add signup to signup table: %v", err)
	}
	return key, nil
}

// IsSignupVerified checks the database to see if an email address exits and is verified.
func IsSignupVerified(ctx context.Context, email string) (bool, error) {
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return false, fmt.Errorf("could not create datastore client: %v", err)
	}

	query := datastore.NewQuery("Signup").
		Filter("email =", email).
		Filter("verified =", true)

	var theSignup Signup
	it := client.Run(ctx, query)
	for {
		_, err := it.Next(&theSignup)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, fmt.Errorf("signup not in DB: %v", err)
		}
		return true, nil
	}
	return false, fmt.Errorf("signup does not exist in DB")
}

// IsCodeAvailable checks the database to see if a code is unique (true)/already used (false)
func IsCodeAvailable(ctx context.Context, code string) (bool, error) {
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return false, fmt.Errorf("could not create datastore client: %v", err)
	}

	query := datastore.NewQuery("Signup").
		Filter("code =", code)

	var theSignup Signup
	it := client.Run(ctx, query)
	for {
		_, err := it.Next(&theSignup)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, fmt.Errorf("signup not in DB: %v", err)
		}
		return false, fmt.Errorf("code not available: %v", err)
	}
	return true, nil
}

// Given a code, return the datastore key and signup
func GetSignup(ctx context.Context, code string) (*datastore.Key, Signup, error) {
	var theSignup Signup
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, Signup{}, fmt.Errorf("could not create datastore client: %v", err)
	}

	query := datastore.NewQuery("Signup").
		Filter("code =", code)

	it := client.Run(ctx, query)
	for {
		k, err := it.Next(&theSignup)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, Signup{}, fmt.Errorf("signup not in DB: %v", err)
		}
		return k, theSignup, nil
	}
	return nil, Signup{}, fmt.Errorf("signup not in DB: %v", err)
}

// MarkVerified marks the signup as verified with the given ID.
func MarkVerified(ctx context.Context, code string) error {
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return fmt.Errorf("could not create datastore client: %v", err)
	}
	k, _, err := GetSignup(ctx, code)
	if err != nil {
		return fmt.Errorf("could not find code in DB: %v", err)
	}
	// Create a key using the given integer ID.
	signup := new(Signup)

	// In a transaction load each signup, set verified to true and store.
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := tx.Get(k, signup); err != nil {
			return errors.New("no such verification code")
		}
		if signup.IsVerified {
			return errors.New("signup already verified")
		}
		signup.IsVerified = true
		_, err := tx.Put(k, signup)
		return err
	})
	return err
}

// Given a code, return the datastore key and signup
func GetSite(ctx context.Context, siteCode string) (Site, error) {
	var theSite Site
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return Site{}, fmt.Errorf("could not create datastore client: %v", err)
	}

	query := datastore.NewQuery("Site").
		Filter("code =", siteCode)

	it := client.Run(ctx, query)
	for {
		_, err := it.Next(&theSite)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return Site{}, fmt.Errorf("site not in DB: %v", err)
		}
		return theSite, nil
	}
	return Site{}, fmt.Errorf("site not in DB: %v", err)
}

// GetSignupCode gets the signup code matching the given email address.
// This should only be called during testing.
func GetSignupCode(ctx context.Context, email string) (string, error) {
	var theSignup Signup
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return "", fmt.Errorf("could not create datastore client: %v", err)
	}

	query := datastore.NewQuery("Signup").
		Filter("email =", email)

	it := client.Run(ctx, query)
	for {
		_, err := it.Next(&theSignup)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("signup not in DB: %v", err)
		}
		return theSignup.VerificationCode, nil
	}
	return "", fmt.Errorf("signup not in DB: %v", err)
}

// AddSite adds a site with the given verification code to the datastore,
// returning the key of the newly created entity.
// Should only be called during testing
func AddSite(ctx context.Context, siteName, siteCode, rootURL string) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("could not create datastore client: %v", err)
	}
	key := datastore.IncompleteKey("Site", nil)
	site := &Site{
		CreationTimestamp: time.Now(),
		SiteName:          siteName,
		Code:              siteCode,
		RootURL:           rootURL,
		VerifiedURL:       rootURL,
	}
	if _, err := client.Put(ctx, key, site); err != nil {
		return nil, fmt.Errorf("could not add site to site table: %v", err)
	}
	return key, nil
}
