# qxb

qxb manages initial access to the web server.
It handles:

* static files
* session management

## Static Files
In development, routes are compared against the filesystem using `os.Stat` on every request.
In production, we cache the filesystem and use a lookup instead of calling `os.Stat`.
(Note that we are not using an embedded filesystem because it doesn't preserve timestamps on the entries.)

## Session Management
We're using the session management from
* https://themsaid.com/building-secure-session-manager-in-go
* https://themsaid.com/csrf-protection-go-web-applications
* https://themsaid.com/session-authentication-go

We're using magic links to authenticate users.
This should be replaced with hashed passwords in the future.