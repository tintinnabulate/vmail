# registration

## Email verification
* Looks like I will have to roll my own

1. Present 'Create Account' form to user
2. On submit, take email address from 'Create Account' user form.
3. Check with [GoValidator](https://github.com/asaskevich/govalidator) for email address format validity.
4. Generate a unique email verification code (e.g. 6 random numbers [0-9]
5. `if` generated code is already marked `InUse`, goto (3), `else` mark as `InUse` and goto (5)
6. Try to `net/smtp.SendMail()` them a unique (hash?) email verification code that expires when session expires.
7. `if` error from `net/smtp.SendMail()`, flash error to user and goto (1), else goto (8)
8. Set an expectation somewhere for that generated verification code
9. Present verification code fill-in form to user
10. On submit, check that filled-in code matches expected code for session
11. Present 'Log in' form to user
