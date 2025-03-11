### **📖 Responders Package – Action-Domain-Responder (ADR) in Go**

#### **Overview**
The **responders** package is part of an **Action-Domain-Responder (ADR)** architecture in Go. This package is responsible for formatting and returning HTTP responses, ensuring a clean separation between request handling, business logic, and response generation.

ADR is an alternative to the traditional **Model-View-Controller (MVC)** pattern, specifically designed for **web applications and APIs**. It provides a clear separation between:
- **Action**: Handles HTTP requests and calls the Domain logic.
- **Domain**: Encapsulates business logic, separate from the web layer.
- **Responder**: Formats and returns the HTTP response (HTML, JSON, XML, etc.).

For a deep dive into ADR, refer to **Paul M. Jones' original article**:  
🔗 [Action-Domain-Responder: A Web-Specific Refinement of MVC](http://pmjones.io/adr/)

---

### **🎯 What Are Responders?**
A **Responder** in ADR is responsible for **formatting and sending HTTP responses**. Unlike MVC’s **views**, which often mix concerns, Responders are dedicated **only** to responding with structured data, making them highly reusable and testable.

#### **Responsibilities of a Responder:**
✔ Determine the appropriate response format (HTML, JSON, XML, etc.).  
✔ Handle success and error responses in a consistent way.  
✔ Ensure a clean separation of response handling from business logic.  
✔ Make it easier to support **HTMX**, REST APIs, and traditional web applications.

---

### **🛠 Implementing Responders in Go**
The **responders** package provides a reusable way to handle responses in a Go web application.

#### **1. Define a Responder Interface**
A generic responder interface helps enforce structure.

```go
package responder

import "net/http"

type Responder interface {
    Respond(w http.ResponseWriter, data any, err error)
}
```

#### **2. Implement a JSON Responder**
A **JSONResponder** returns structured JSON responses.

```go
package responder

import (
    "encoding/json"
    "net/http"
)

type JSONResponder struct{}

func (r *JSONResponder) Respond(w http.ResponseWriter, data any, err error) {
    w.Header().Set("Content-Type", "application/json")

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(data)
}
```

#### **3. Implement an HTMX Responder**
For **HTMX**, we return HTML fragments instead of JSON.

```go
package responder

import (
    "html/template"
    "net/http"
)

type HTMXResponder struct {
    Tmpl *template.Template
}

func (r *HTMXResponder) Respond(w http.ResponseWriter, data any, err error) {
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    r.Tmpl.ExecuteTemplate(w, "user-row.html", data)
}
```

---

### **🚀 Using Responders in an ADR Web Application**
#### **1. Inject the Responder into an Action**
In an **ADR-based Go application**, the **Action** calls the **Domain** logic and then delegates response handling to the **Responder**.

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

#### **2. Register Routes with Different Responders**
Depending on the request type, you can return **JSON, HTML, or other formats**.

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

    // JSON Responder (for APIs)
    jsonResponder := &responder.JSONResponder{}
    r.Post("/users/json", (&action.CreateUserAction{Service: userService, Responder: jsonResponder}).ServeHTTP)

    // HTMX Responder (for partial HTML updates)
    htmxResponder := &responder.HTMXResponder{Tmpl: tmpl}
    r.Post("/users", (&action.CreateUserAction{Service: userService, Responder: htmxResponder}).ServeHTTP)

    http.ListenAndServe(":8080", r)
}
```

---

### **💡 Benefits of Using ADR and Responders**
✅ **Separation of Concerns** → The **Action**, **Domain**, and **Responder** are independent.  
✅ **Supports APIs and HTMX** → Use the same **Action** with different **Responders**.  
✅ **Flexible Response Handling** → Return **JSON, HTML, XML, or other formats** easily.  
✅ **Testability** → You can unit test **Responders** without setting up an HTTP server.

---

### **📚 Additional Reading**
🔗 **Paul M. Jones' ADR Article**: [http://pmjones.io/adr/](http://pmjones.io/adr/)  
🔗 **GitHub - ADR Example Repo**: [https://github.com/pmjones/adr-example](https://github.com/pmjones/adr-example)  
🔗 **HTMX Documentation**: [https://htmx.org/](https://htmx.org/)

---

### **🎯 Summary**
The **responders** package helps implement the **Responder** portion of the **ADR** pattern in Go. It provides reusable components for handling **JSON, HTML (HTMX), and other response types**, making it easy to support both **REST APIs and modern web applications**.
