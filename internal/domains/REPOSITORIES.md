The separation of **business logic** and **data access** is an important consideration in ADR. Let’s clarify where repositories belong.

---

## **Where Should Repositories Go?**
Repositories **do not** belong in the **Domain** layer if you want a strict separation of business logic from data access. Instead, repositories should be placed in a **separate package**, commonly named **storage**, **repository**, or **persistence**.

### **Revised Layer Breakdown**
| Layer                    | Responsibility                                                   | Examples                         |
|--------------------------|------------------------------------------------------------------|----------------------------------|
| **Action**               | Handles HTTP requests, extracts parameters, calls the **Domain** | `CreateUserAction`               |
| **Domain**               | Pure business logic, no knowledge of data storage                | `UserService`                    |
| **Repository (Storage)** | Handles database interactions, abstracted from the **Domain**    | `UserRepository`                 |
| **Responder**            | Formats responses (JSON, HTML, etc.)                             | `JSONResponder`, `HTMXResponder` |

---

## **Revised Implementation – Moving Repositories to a Storage Package**
### **1. Define a Repository Interface in `domains`**
Repositories are **abstracted** inside the **domain** layer but implemented outside it.

```go
package domains

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// UserRepository defines an interface for storage operations
type UserRepository interface {
    Save(user User) error
    FindByID(id string) (User, error)
}
```

The **UserRepository** interface stays in `domains`, ensuring that the **domain layer does not depend on storage details**.

---

### **2. Implement the Repository in a Separate `storage` Package**
The actual **data storage logic** is implemented in a separate `storage` package.

```go
package storage

import (
    "errors"
    "myapp/domains"
)

type InMemoryUserRepo struct {
    data map[string]domains.User
}

// Constructor function
func NewInMemoryUserRepo() *InMemoryUserRepo {
    return &InMemoryUserRepo{data: make(map[string]domains.User)}
}

// Implements UserRepository interface
func (repo *InMemoryUserRepo) Save(user domains.User) error {
    repo.data[user.ID] = user
    return nil
}

func (repo *InMemoryUserRepo) FindByID(id string) (domains.User, error) {
    user, exists := repo.data[id]
    if !exists {
        return domains.User{}, errors.New("user not found")
    }
    return user, nil
}
```

This ensures that **storage logic is separate** and can be easily replaced (e.g., replacing an in-memory store with a SQL database).

---

### **3. Keep the Domain Layer Pure**
The **domain layer** should only call **repository interfaces**, without knowing their implementation details.

```go
package domains

import (
    "errors"
    "github.com/google/uuid"
)

type UserService struct {
    Repo UserRepository // Dependency injection
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

This structure ensures that the **Domain layer remains pure** and **storage concerns are handled externally**.

---

## **Updated `main.go` – Injecting Dependencies**
With this revised structure, we **inject the repository implementation into the domain service** at runtime.

```go
package main

import (
    "html/template"
    "net/http"

    "github.com/go-chi/chi/v5"
    "myapp/actions"
    "myapp/domains"
    "myapp/responders"
    "myapp/storage"
)

func main() {
    r := chi.NewRouter()

    tmpl := template.Must(template.ParseFiles("templates/user-row.html"))

    userRepo := storage.NewInMemoryUserRepo() // Storage layer
    userService := &domains.UserService{Repo: userRepo} // Business logic layer

    // JSON API
    jsonResponder := &responders.JSONResponder{}
    r.Post("/users/json", (&actions.CreateUserAction{Service: userService, Responder: jsonResponder}).ServeHTTP)

    // HTMX UI updates
    htmxResponder := &responders.HTMXResponder{Tmpl: tmpl}
    r.Post("/users", (&actions.CreateUserAction{Service: userService, Responder: htmxResponder}).ServeHTTP)

    http.ListenAndServe(":8080", r)
}
```

---

## **🚀 Why This Structure is Better**
✅ **Strict Separation of Concerns** → The domain layer **doesn’t know about databases**.  
✅ **Testability** → You can mock repositories for testing the domain logic.  
✅ **Flexibility** → Swap **storage implementations** (e.g., switch from in-memory to PostgreSQL).  
✅ **Encapsulation** → The **domain layer is reusable** across different transports (CLI, jobs, APIs).

---

## **📚 Updated Summary for the `domains` README**
### **Where Does the Domain Layer Fit?**
The **Domain** layer contains:

✔ **Business rules** that enforce how data should be processed.  
✔ **Interfaces** that define required behaviors (e.g., `UserRepository`).  
✔ **Services** that act as the API for business logic.

The **Domain layer does not:**

❌ Contain direct database queries (handled by `storage`).  
❌ Know about HTTP, JSON, or response formatting (handled by `action` and `responder`).

---

### **📚 Additional Reading**
🔗 **Paul M. Jones' ADR Article**: [http://pmjones.io/adr/](http://pmjones.io/adr/)  
🔗 **Repository Pattern in Go**: [https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091](https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091)  
🔗 **Clean Architecture in Go**: [https://8thlight.com/insights/guides/clean-architecture/](https://8thlight.com/insights/guides/clean-architecture/)
