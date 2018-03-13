% Registration ICYfication
% Implementation Thoughts - Issue 0.1, Draft
% \today{}
\newpage{}

# Implementation thoughts

## Testing

This is Test-First Development, so I'll be writing tests before code.

* [Testing](https://golang.org/pkg/testing/) - go's builtin testing package, for unit-testing.
* [Quick](https://golang.org/pkg/testing/quick/) - go's builtin QuickCheck implementation, for unit-testing.
* [Table-driven tests](https://github.com/golang/go/wiki/TableDrivenTests) - example of this: [fmt_test.go](https://golang.org/src/fmt/fmt_test.go), for Boundary-Value analysis, and saving writing lines of code.
* [Frisby](https://github.com/verdverm/frisby) for application-level testing, i.e. Frisby visits the site, tries to register, pay, etc.
* [stripe-mock](https://github.com/stripe/stripe-mock) - for testing mock Stripe API without actually hitting the real API.

## Web framework

* [Echo](https://echo.labstack.com/guide)

## Email verification
* Looks like I will have to roll my own
* Should this be a small, self-contained go microservice? It could work with it's own db, or just a db table. It lends itself nicely to a microservice.

1. Present 'Create Account' form to user,
2. On submit, take email address from 'Create Account' user form,
3. Check with [GoValidator](https://github.com/asaskevich/govalidator) for email address format validity,
4. Generate a unique email verification code (e.g. 6 random numbers [0-9],
5. `if` generated code is already marked `InUse`, goto (4), `else` mark as `InUse` and goto (6),
6. Try to [SendMail()](https://golang.org/pkg/net/smtp/#SendMail) them a unique (hash?) email verification code that expires when session expires,
7. `if` error from [SendMail()](https://golang.org/pkg/net/smtp/#SendMail), flash error to user and goto (1), `else` goto (8),
8. Set an expectation somewhere for that generated verification code,
9. Present verification code fill-in form to user,
10. On submit, check that filled-in code matches expected code for session,
11. Present 'Log in' form to user,
12. End.

## Taking payment

* [stripe-go](https://stripe.com/docs/checkout/go) - Just use Stripe :)

# Change history

Date | Issue | Note
---|---|---
2017-03-12 | 0.1, Draft | Initial draft |
