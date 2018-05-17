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
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
