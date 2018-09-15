<a href="https://goreportcard.com/report/github.com/tintinnabulate/registration"><img src="https://goreportcard.com/badge/github.com/tintinnabulate/registration" /></a>

# Signup

Really this should be called 'Email verifier'...

## Endpoints

### `POST /signup/{email_address}`

1. Generate random verification code `XXXXXX`
2. Check that code does not already exist in database. If it does exist, GOTO (1), else GOTO (3).
3. Send email to `email_address` containing link to `/verify/XXXXXX`
4. Store in database the row `creation_date: time.Now(), email: email_address, code: XXXXXX, is_verified: false`
5. Response JSON: `{"Address": "email_address", "Success": true, "Note": ""}`

### `GET /signup/{email_address}`

1. Check `email_address` in database and `verified = true`
2. Response JSON: `{"Address": "email_address", "Success": true, "Note": ""}`

### `GET /verify/{code}`

1. Look up `code = XXXXXX` in database. If `XXXXXX` exists, GOTO 2, else GOTO 3.
2. If `is_verified: false` GOTO 4, else GOTO 5
3. Response JSON: `{"Code": "XXXXXX", "Success": false, "Note": "no such verification code"}`
4. Mark `is_verified` as `true` in database, GOTO 6
5. Response JSON: `{"Code": "XXXXXX", "Success": false, "Note": "signup already verified"}`
6. Response JSON: `{"Code": "XXXXXX", "Success": true, "Note": ""}`

## Notes

I found this incredibly useful in helping me design testable code using Google App Engine: 
https://www.compoundtheory.com/testing-go-http-handlers-in-google-app-engine-with-mux-and-higher-order-functions/
