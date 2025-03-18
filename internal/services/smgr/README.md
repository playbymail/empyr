# Session Manager

Uses code from

* https://themsaid.com/building-secure-session-manager-in-go
* https://themsaid.com/csrf-protection-go-web-applications

## Session Based Authentication

From https://themsaid.com/session-authentication-go

The idea behind session-based authentication is:

    Client                                    Server
    |                                          |
    |  1. User visits the web app              |
    |----------------------------------------> |
    |                                          |
    |  2. Server creates a session and sends   |
    |     its ID back                          |
    | <--------------------------------------- |
    |                                          |
    |  3. User submits login form              |
    |----------------------------------------> |
    |                                          |
    |  4. Server verifies credentials and      |
    |     stores the user identifier in the    |
    |     session                              |
    |                                          |
    |  5. Server sends a success response      |
    | <----------------------------------------|
    |                                          |
    |  6. Client makes authenticated requests  |
    |     (Cookie is automatically sent)       |
    |----------------------------------------> |
    |                                          |
    |  7. Server retrieves session data and    |
    |     extracts the user identifier         |
    |                                          |
    |  8. Server processes request and         |
    |     responds                             |
    | <----------------------------------------|

