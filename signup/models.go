/*
	Implementation Note:
		None.
	Filename:
		models.go
*/

package main

// Email holds our JSON response for GET and POST /signup/{email}
type Email struct {
	Address string `json:"address"`
	Success bool   `json:"success"`
	Note    string `json:"note"`
}

// Verification holds our JSON response for GET /verify/{code}
type Verification struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Note    string `json:"note"`
}
