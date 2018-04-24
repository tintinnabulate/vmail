// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START all]

// A simple command-line task list manager to demonstrate using the
// cloud.google.com/go/datastore package.
package main

import (
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
	signup := new(Signup)
	signup.CreationTimestamp = time.Now()
	signup.Email = email
	signup.VerificationCode = code
	signup.IsVerified = false
	foo, err := datastore.Put(ctx, key, signup)
	return foo, err
}

// MarkDone marks the signup as verified with the given ID.
func MarkVerified(r *http.Request, code string) error {
	ctx := appengine.NewContext(r)
	// Create a key using the given integer ID.
	key := datastore.NewKey(ctx, "Signup", code, 0, nil)

	// In a transaction load each task, set done to true and store.
	err := datastore.RunInTransaction(ctx, func(tx context.Context) error {
		var signup Signup
		if err := datastore.Get(tx, key, &signup); err != nil {
			return err
		}
		signup.IsVerified = true
		_, err := datastore.Put(tx, key, &signup)
		return err
	}, nil)
	return err
}

/*

// ListSignups returns all the tasks in ascending order of creation time.
func ListSignups(r *http.Request) ([]*Signup, error) {
	var signups []*Signup

	ctx := appengine.NewContext(r)
	client, err := datastore.NewClient(ctx, config.ProjectID)
	checkErr(err)

	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery("Signup").Order("created")
	keys, err := client.GetAll(ctx, query, &signups)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		signups[i].id = key.ID
	}

	return signups, nil
}

// DeleteSignup deletes the task with the given ID.
func DeleteSignup(r *http.Request, signupID int64) error {
	ctx := appengine.NewContext(r)
	client, err := datastore.NewClient(ctx, config.ProjectID)
	checkErr(err)
	return client.Delete(ctx, datastore.IDKey("Signup", signupID, nil))
}
*/
