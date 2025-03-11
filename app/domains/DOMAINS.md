# **📖 Domains Package – Action-Domain-Responder (ADR) in Go**

## **Overview**
The **domains** package is responsible for encapsulating **business logic** in an **Action-Domain-Responder (ADR)** structured Go application. It serves as the core of the application, independent of the **HTTP layer**, ensuring that all rules and logic are reusable across different interfaces (e.g., web, CLI, background jobs).

ADR is a web-specific alternative to **Model-View-Controller (MVC)** that clearly separates concerns:
- **Action** → Receives HTTP requests, extracts input, and calls the **Domain**.
- **Domain** → Contains business logic and interacts with data storage.
- **Responder** → Formats and returns the response (JSON, HTML, XML, etc.).

🔗 **For a deep dive into ADR, read Paul M. Jones' original article:**  
[Action-Domain-Responder: A Web-Specific Refinement of MVC](http://pmjones.io/adr/)

---

## **🎯 What is the Domain Layer?**
In ADR, the **Domain** layer is where the actual application logic resides. It:
✔ Encapsulates business rules independent of HTTP or any transport mechanism.  
✔ Interacts with databases, external APIs, and services.  
✔ Provides a **clean API** for the **Action** layer to call.  
✔ Keeps business logic separate from request handling (Action) and response formatting (Responder).

Unlike MVC’s **models**, which often mix database operations and business logic, the Domain layer in ADR **strictly** focuses on **business operations**.

---

## **🛠 Implementing the Domain Layer in Go**
The **domains** package provides services and repositories for handling core application logic.

### **1. Define a Domain Entity (User)**
Entities represent business objects and contain logic related to their domain.

```go
package domain

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
```

---

### **2. Define a Repository Interface**
Repositories abstract data storage, making the domain logic independent of the database.

```go
package domain

type UserRepository interface {
    Save(user User) error
    FindByID(id string) (User, error)
}
```

---

### **3. Implement an In-Memory User Repository**
A basic repository for testing purposes.

```go
package domain

import "errors"

type InMemoryUserRepo struct {
    data map[string]User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
    return &InMemoryUserRepo{data: make(map[string]User)}
}

func (repo *InMemoryUserRepo) Save(user User) error {
    repo.data[user.ID] = user
    return nil
}

func (repo *InMemoryUserRepo) FindByID(id string) (User, error) {
    user, exists := repo.data[id]
    if !exists {
        return User{}, errors.New("user not found")
    }
    return user, nil
}
```

---

### **4. Implement a Domain Service (UserService)**
Services handle business logic, calling repositories when needed.

```go
package domain

import (
    "errors"
    "github.com/google/uuid"
)

type UserService struct {
    Repo UserRepository
}

func (s *UserService) CreateUser(username, email string) (User, error) {
    if username == "" || email == "" {
        return User{}, errors.New("invalid input")
    }

    user := User{
        ID:       uuid.New().String(),
        Username: username,
        Email:    email,
    }

    err := s.Repo.Save(user)
    return user, err
}

func (s *UserService) GetUser(id string) (User, error) {
    return s.Repo.FindByID(id)
}
```

---

## **🚀 Using the Domain Layer in an ADR Web Application**
### **1. Injecting the Domain into an Action**
Actions call **Domain services** instead of dealing with data storage directly.

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

---

### **2. Registering Routes with Domain Services**
In the main function, inject the **domain service** and use it in **Actions**.

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

    userRepo := domain.NewInMemoryUserRepo()
    userService := &domain.UserService{Repo: userRepo}

    // JSON Responder (API)
    jsonResponder := &responder.JSONResponder{}
    r.Post("/users/json", (&action.CreateUserAction{Service: userService, Responder: jsonResponder}).ServeHTTP)

    // HTMX Responder (UI Updates)
    htmxResponder := &responder.HTMXResponder{Tmpl: tmpl}
    r.Post("/users", (&action.CreateUserAction{Service: userService, Responder: htmxResponder}).ServeHTTP)

    http.ListenAndServe(":8080", r)
}
```

---

## **💡 Benefits of Using a Domain Layer in ADR**
✅ **Separation of Concerns** → The **Domain** is independent of HTTP and UI layers.  
✅ **Testable** → Business logic can be tested without setting up an HTTP server.  
✅ **Reusability** → The same Domain logic can be used in CLI tools, background jobs, or APIs.  
✅ **Encapsulated Business Logic** → Keeps Actions thin and focused on request handling.

---

## **📚 Additional Reading**
🔗 **Paul M. Jones' ADR Article**: [http://pmjones.io/adr/](http://pmjones.io/adr/)  
🔗 **GitHub - ADR Example Repo**: [https://github.com/pmjones/adr-example](https://github.com/pmjones/adr-example)  
🔗 **Google Clean Architecture (Similar Concepts)**: [https://8thlight.com/insights/guides/clean-architecture/](https://8thlight.com/insights/guides/clean-architecture/)

---

## **🎯 Summary**
The **domains** package provides the **business logic** in an ADR-based Go application. By structuring logic **independent of HTTP**, we create a **testable, scalable, and maintainable** architecture that supports both **API-driven and UI-based** web applications.