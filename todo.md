# TODO

* Use gRPC, not this horrendous 'REST-ful' 'API'.
* Check out https://gokit.io/ also.
* Read https://12factor.net/
* Require OAUTH2?
    * [https://cloud.google.com/go/getting-started/authenticate-users]()
    * [https://medium.com/@hfogelberg/the-black-magic-of-oauth-in-golang-part-1-3cef05c28dde]()
    * [https://cloud.google.com/appengine/docs/standard/go/users/]()
    * [https://github.com/google/google-api-go-client/blob/master/GettingStarted.md]()
    * Also see my starred repos
* Don't use `ErrCheck(err)` everywhere...
	* At the moment I mostly use `ErrCheck`, a function that calls `log.Fatal` if `err != nil`. This is laaaazy. I need to handle the error in an appropriate way for that particular context.
