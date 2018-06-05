# Signup

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

.

## TODO

* Require OAUTH2
    * [https://cloud.google.com/go/getting-started/authenticate-users]()
    * [https://medium.com/@hfogelberg/the-black-magic-of-oauth-in-golang-part-1-3cef05c28dde]()
    * [https://cloud.google.com/appengine/docs/standard/go/users/]()
    * [https://github.com/google/google-api-go-client/blob/master/GettingStarted.md]()
    * Also see my starred repos
* Don't use `ErrCheck(err)` everywhere...
	* At the moment I mostly use `ErrCheck`, a function that calls `log.Fatal` if `err != nil`. This is laaaazy. I need to handle the error in an appropriate way for that particular context.
