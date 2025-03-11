# **📖 Actions Package – Action-Domain-Responder (ADR) in Go**

## **Overview**
The **actions** package is responsible for handling **HTTP requests** in an **Action-Domain-Responder (ADR)** structured Go web application. Actions serve as the entry point for handling web requests, delegating business logic to the **Domain** and response formatting to the **Responder**.

ADR is a web-specific alternative to **Model-View-Controller (MVC)** that clearly separates concerns:
- **Action** → Receives HTTP requests, extracts input, and calls the Domain.
- **Domain** → Encapsulates the business logic, separate from HTTP concerns.
- **Responder** → Formats and returns the response (JSON, HTML, XML, etc.).

🔗 **For a deep dive into ADR, read Paul M. Jones' original article:**  
[Action-Domain-Responder: A Web-Specific Refinement of MVC](http://pmjones.io/adr/)

---

## **🎯 What Are Actions?**
In the ADR pattern, an **Action** acts as an HTTP request handler. It:
✔ Extracts and validates input parameters from the request.  
✔ Calls the **Domain** logic to process the request.  
✔ Delegates response handling to the **Responder**.  
✔ Does **not** contain business logic (which stays in the **Domain**).

Unlike MVC’s **controllers**, Actions do **only one thing per request**, making them highly modular and easy to test.

---

## **🛠 Implementing Actions in Go**
The **actions** package provides reusable components to handle HTTP requests cleanly in an ADR web application.

### **1. Define a Generic Action Interface**
A reusable interface for all actions ensures consistency.

```go
package action

import "net/http"

type Action interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}
```

### **2. Implement a CreateUserAction**
A **CreateUserAction** extracts request data, calls the **Domain**, and passes results to the **Responder**.

```go
package action

import (
    "net/http"
    "myapp/domain"
    "myapp/responder"
)

type CreateUserAction struct {
    Service   *domain.UserService
    Responder responder.Responder
}

func (a *CreateUserAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    email := r.FormValue("email")

    user, err := a.Service.CreateUser(username, email)
    a.Responder.Respond(w, user, err)
}
```

### **3. Implement an UpdateUserAction**
For updating users, we follow the same structure.

```go
package action

import (
    "net/http"
    "myapp/domain"
    "myapp/responder"
)

type UpdateUserAction struct {
    Service   *domain.UserService
    Responder responder.Responder
}

func (a *UpdateUserAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    username := r.FormValue("username")
    email := r.FormValue("email")

    user, err := a.Service.UpdateUser(id, username, email)
    a.Responder.Respond(w, user, err)
}
```

---

## **🚀 Using Actions in an ADR Web Application**
### **1. Registering Actions with a Router**
Actions are mapped to HTTP routes, allowing them to process requests.

```go
package main

import (
    "html/template"
    "net/http"

    "github.com/go-chi/chi/v5"
    "myapp/action"
    "myapp/domain"
    "myapp/responder"
)

func main() {
    r := chi.NewRouter()

    tmpl := template.Must(template.ParseFiles("templates/user-row.html"))

    userRepo := &domain.InMemoryUserRepo{}
    userService := &domain.UserService{Repo: userRepo}

    // JSON Responder for API requests
    jsonResponder := &responder.JSONResponder{}
    createUserAction := &action.CreateUserAction{Service: userService, Responder: jsonResponder}
    r.Post("/users/json", createUserAction.ServeHTTP)

    // HTMX Responder for dynamic UI updates
    htmxResponder := &responder.HTMXResponder{Tmpl: tmpl}
    r.Post("/users", (&action.CreateUserAction{Service: userService, Responder: htmxResponder}).ServeHTTP)

    // Start the server
    http.ListenAndServe(":8080", r)
}
```

---

## **💡 Benefits of Using Actions in ADR**
✅ **Separation of Concerns** → Business logic is in the **Domain**, not in the Action.  
✅ **Single Responsibility** → Each Action handles **only one request type**.  
✅ **Reusability** → The same Action can work with different Responders.  
✅ **Scalability** → Adding new endpoints is simple and keeps the codebase clean.

---

## **📚 Additional Reading**
🔗 **Paul M. Jones' ADR Article**: [http://pmjones.io/adr/](http://pmjones.io/adr/)  
🔗 **GitHub - ADR Example Repo**: [https://github.com/pmjones/adr-example](https://github.com/pmjones/adr-example)  
🔗 **HTMX Documentation**: [https://htmx.org/](https://htmx.org/)

---

## **🎯 Summary**
The **actions** package provides a structured approach to handling HTTP requests in an ADR-based Go application. By following this pattern, we achieve **cleaner, more testable, and maintainable** web applications while supporting both **HTMX-powered UI updates and API responses**.