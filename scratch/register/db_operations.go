/*
	Implementation Note:
		None.

	Filename:
		db_operations.go
*/

package main

import (
	"fmt"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

// AddSignup adds a signup with the given verification code to the datastore,
// returning the key of the newly created entity.
func StashRegistrationForm(ctx context.Context, regform *RegistrationForm) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "RegistrationForm", regform.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, regform)
	return k, err
}

// GetSignupCode gets the signup code matching the given email address.
// This should only be called during testing.
func GetRegistrationForm(ctx context.Context, email string) (RegistrationForm, error) {
	q := datastore.NewQuery("RegistrationForm").Filter("Email_Address =", email)
	var regforms []RegistrationForm
	if _, err := q.GetAll(ctx, &regforms); err != nil {
		return RegistrationForm{}, err
	}
	if len(regforms) < 1 {
		return RegistrationForm{}, fmt.Errorf("Email not in database: %s", email)
	}
	return regforms[0], nil
}

func AddUser(ctx context.Context, user *User) (*datastore.Key, error) {
	key := datastore.NewKey(ctx, "User", user.Email_Address, 0, nil)
	k, err := datastore.Put(ctx, key, user)
	return k, err
}
