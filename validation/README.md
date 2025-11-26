# Validation Package

A simple, flexible validation library for Go with struct tags, map validation, and custom rules support.

## Table of Contents

- [What is Validation?](#what-is-validation)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Validation Methods](#validation-methods)
- [Available Rules](#available-rules)
- [Error Handling](#error-handling)
- [Complete Examples](#complete-examples)
- [Custom Rules](#custom-rules)
- [Best Practices](#best-practices)

---

## What is Validation?

The validation package helps you **check if data is correct** before using it in your application. It's like a security guard that checks if everything is in order before letting data through.

**Simple Analogy:**
- **Validation** = Security guard checking IDs at the entrance
- **Rules** = Requirements (age must be 18+, email must be valid, etc.)
- **Errors** = List of problems found
- **Struct Tags** = Instructions written on the form
- **Custom Rules** = Your own special requirements

---

## Features

- ✅ **Struct Tag Validation** - Validate using `validate:""` tags
- ✅ **Map Validation** - Validate map[string]any data
- ✅ **Single Field Validation** - Validate individual fields
- ✅ **20+ Built-in Rules** - Email, phone, password, URL, etc.
- ✅ **Custom Rules** - Create your own validation logic
- ✅ **Multiple Errors** - Get all validation errors at once
- ✅ **Type-Safe** - Works with Go types
- ✅ **Zero Dependencies** - Only uses standard library (except UUID)
- ✅ **Easy to Use** - Simple, intuitive API

---

## Installation

```bash
go get github.com/fatkulnurk/foundation/validation
```

**Dependencies:**
- Go 1.25 or higher
- github.com/google/uuid (for UUID validation)

---

## Quick Start

### 1. Struct Validation

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

type User struct {
    Name  string `json:"name" validate:"required,strminlen=3,strmaxlen=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"nummin=18,nummax=100"`
}

func main() {
    user := User{
        Name:  "Jo",  // Too short
        Email: "invalid-email",  // Invalid format
        Age:   15,   // Too young
    }
    
    errs := validation.ValidateStruct(user)
    
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    }
}
```

**Output:**
```
name: must be at least 3 characters
email: format email tidak valid
age: must be at least 18.00
```

### 2. Map Validation

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

func main() {
    data := map[string]any{
        "username": "john",
        "password": "weak",
        "age":      15,
    }
    
    rules := map[string][]validation.Rule{
        "username": {
            validation.Required(""),
            validation.StrMinLength(6, ""),
        },
        "password": {
            validation.Required(""),
            validation.Password(""),
        },
        "age": {
            validation.NumMin(18, ""),
        },
    }
    
    errs := validation.ValidateMap(data, rules)
    
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    }
}
```

---

## Validation Methods

### 1. ValidateStruct

Validates a struct using `validate` tags.

```go
func ValidateStruct(s any) Errors
```

**Usage:**
```go
type Product struct {
    Name  string `json:"name" validate:"required,strminlen=3"`
    Price int    `json:"price" validate:"required,nummin=0"`
}

product := Product{Name: "AB", Price: -10}
errs := validation.ValidateStruct(product)
```

**Tag Format:**
- Use `validate` tag to specify rules
- Separate multiple rules with commas
- Field name comes from `json` tag (or struct field name)

### 2. ValidateMap

Validates a map with programmatic rules.

```go
func ValidateMap(data map[string]any, rules map[string][]Rule) Errors
```

**Usage:**
```go
data := map[string]any{
    "email": "test@example.com",
    "age":   25,
}

rules := map[string][]validation.Rule{
    "email": {validation.Required(""), validation.Email("")},
    "age":   {validation.NumMin(18, "")},
}

errs := validation.ValidateMap(data, rules)
```

### 3. Validate (Single Field)

Validates a single field with a rule string.

```go
func Validate(field string, value any, rule string) *Error
```

**Usage:**
```go
err := validation.Validate("email", "test@example.com", "required,email")
if err != nil {
    fmt.Println(err.Message)
}
```

---

## Available Rules

### String Rules

#### required
Value cannot be empty or whitespace.

```go
// Struct tag
`validate:"required"`

// Programmatic
validation.Required("")
```

#### strminlen=N
Minimum string length.

```go
// Struct tag
`validate:"strminlen=3"`

// Programmatic
validation.StrMinLength(3, "")
```

#### strmaxlen=N
Maximum string length.

```go
// Struct tag
`validate:"strmaxlen=50"`

// Programmatic
validation.StrMaxLength(50, "")
```

#### email
Valid email format (must contain @).

```go
// Struct tag
`validate:"email"`

// Programmatic (use tag parsing)
validation.Validate("email", value, "email")
```

#### username
Valid username: 6-16 characters, letters, numbers, underscores only.

```go
// Struct tag
`validate:"username"`

// Programmatic
validation.Username("")
```

#### password
Strong password: 8-16 characters, must include uppercase, lowercase, number, and special character.

```go
// Struct tag
`validate:"password"`

// Programmatic
validation.Password("")
```

#### phone
Valid phone number: 10-15 digits, allows spaces, dashes, parentheses, and leading +.

```go
// Struct tag
`validate:"phone"`

// Programmatic
validation.Phone("")
```

#### url
Valid URL with http or https scheme.

```go
// Struct tag
`validate:"url"`

// Programmatic
validation.URLFormat("")
```

#### alphanumeric
Only letters and numbers allowed.

```go
// Struct tag
`validate:"alphanumeric"`

// Programmatic
validation.AlphaNumeric("")
```

### Number Rules

#### nummin=N
Minimum numeric value.

```go
// Struct tag
`validate:"nummin=18"`

// Programmatic
validation.NumMin(18, "")
```

#### nummax=N
Maximum numeric value.

```go
// Struct tag
`validate:"nummax=100"`

// Programmatic
validation.NumMax(100, "")
```

**Supported types:** int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64

### Format Rules

#### date
Valid date in YYYY-MM-DD format.

```go
// Struct tag
`validate:"date"`

// Programmatic
validation.Date("")
```

#### uuid
Valid UUID format.

```go
// Struct tag
`validate:"uuid"`

// Programmatic
validation.Uuid("")
```

#### json
Valid JSON string.

```go
// Struct tag
`validate:"json"`

// Programmatic
validation.Json("")
```

#### hexcolor
Valid hex color code (#RRGGBB).

```go
// Struct tag
`validate:"hexcolor"`

// Programmatic
validation.HexColor("")
```

#### creditcard
Valid credit card number (Luhn algorithm).

```go
// Struct tag
`validate:"creditcard"`

// Programmatic
validation.CreditCard("")
```

#### postalcode
Valid postal code (5 digits for Indonesia).

```go
// Struct tag
`validate:"postalcode"`

// Programmatic
validation.PostalCode("")
```

#### base64
Valid base64 encoded string.

```go
// Struct tag
`validate:"base64"`

// Programmatic
validation.Base64("")
```

### Network Rules

#### ip
Valid IP address (IPv4 or IPv6).

```go
// Struct tag
`validate:"ip"`

// Programmatic (use tag parsing)
validation.Validate("ip", value, "ip")
```

#### ipv4
Valid IPv4 address.

```go
// Struct tag
`validate:"ipv4"`

// Programmatic (use tag parsing)
validation.Validate("ip", value, "ipv4")
```

#### ipv6
Valid IPv6 address.

```go
// Struct tag
`validate:"ipv6"`

// Programmatic (use tag parsing)
validation.Validate("ip", value, "ipv6")
```

---

## Error Handling

### Error Structure

```go
type Error struct {
    Field   string  // Field name
    Message string  // Error message
}

type Errors []Error
```

### Error Methods

#### HasErrors()
Check if there are any errors.

```go
errs := validation.ValidateStruct(user)
if errs.HasErrors() {
    // Handle errors
}
```

#### Error()
Get all errors as a single string.

```go
errs := validation.ValidateStruct(user)
fmt.Println(errs.Error())
// Output: "name: can't be empty; email: is not a valid email"
```

#### ForField(field)
Get errors for a specific field.

```go
errs := validation.ValidateStruct(user)
emailErrors := errs.ForField("email")
for _, err := range emailErrors {
    fmt.Println(err.Message)
}
```

---

## Complete Examples

### Example 1: User Registration

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

type RegisterRequest struct {
    Username string `json:"username" validate:"required,username"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
    Phone    string `json:"phone" validate:"required,phone"`
    Age      int    `json:"age" validate:"required,nummin=18,nummax=100"`
}

func main() {
    req := RegisterRequest{
        Username: "john",  // Too short (min 6)
        Email:    "invalid",  // Invalid email
        Password: "weak",  // Too weak
        Phone:    "123",  // Invalid phone
        Age:      15,  // Too young
    }
    
    errs := validation.ValidateStruct(req)
    
    if errs.HasErrors() {
        fmt.Println("Registration failed:")
        for _, err := range errs {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    } else {
        fmt.Println("Registration successful!")
    }
}
```

### Example 2: Product Validation

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

type Product struct {
    Name        string  `json:"name" validate:"required,strminlen=3,strmaxlen=100"`
    Description string  `json:"description" validate:"strmaxlen=500"`
    Price       float64 `json:"price" validate:"required,nummin=0"`
    Stock       int     `json:"stock" validate:"required,nummin=0"`
    SKU         string  `json:"sku" validate:"required,alphanumeric"`
}

func main() {
    product := Product{
        Name:        "AB",  // Too short
        Description: "A great product",
        Price:       -10,  // Negative price
        Stock:       -5,  // Negative stock
        SKU:         "ABC-123",  // Contains dash (not alphanumeric)
    }
    
    errs := validation.ValidateStruct(product)
    
    if errs.HasErrors() {
        fmt.Println("Product validation failed:")
        for _, err := range errs {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Example 3: API Request Validation

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/fatkulnurk/foundation/validation"
)

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,strminlen=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    
    // Parse JSON
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Validate
    errs := validation.ValidateStruct(req)
    if errs.HasErrors() {
        w.Header().Set("Content-Type", "application/json")
        w.WriteStatus(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]any{
            "error":  "Validation failed",
            "fields": errs,
        })
        return
    }
    
    // Process valid request
    fmt.Fprintf(w, "User created successfully")
}
```

### Example 4: Map Validation with Dynamic Data

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

func main() {
    // Data from form or API
    data := map[string]any{
        "title":   "Hello",
        "content": "This is a blog post",
        "author":  "john_doe",
        "views":   -10,  // Invalid
    }
    
    // Define rules
    rules := map[string][]validation.Rule{
        "title": {
            validation.Required("Title is required"),
            validation.StrMinLength(5, "Title must be at least 5 characters"),
            validation.StrMaxLength(100, "Title too long"),
        },
        "content": {
            validation.Required("Content is required"),
            validation.StrMinLength(10, "Content too short"),
        },
        "author": {
            validation.Required("Author is required"),
            validation.Username("Invalid username format"),
        },
        "views": {
            validation.NumMin(0, "Views cannot be negative"),
        },
    }
    
    errs := validation.ValidateMap(data, rules)
    
    if errs.HasErrors() {
        fmt.Println("Validation errors:")
        for _, err := range errs {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Example 5: Nested Struct Validation

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

type Address struct {
    Street     string `json:"street" validate:"required"`
    City       string `json:"city" validate:"required"`
    PostalCode string `json:"postal_code" validate:"required,postalcode"`
}

type Customer struct {
    Name    string  `json:"name" validate:"required,strminlen=3"`
    Email   string  `json:"email" validate:"required,email"`
    Phone   string  `json:"phone" validate:"required,phone"`
    Address Address `json:"address"`
}

func main() {
    customer := Customer{
        Name:  "John Doe",
        Email: "john@example.com",
        Phone: "+1234567890",
        Address: Address{
            Street:     "123 Main St",
            City:       "New York",
            PostalCode: "12345",
        },
    }
    
    // Validate customer
    errs := validation.ValidateStruct(customer)
    if errs.HasErrors() {
        fmt.Println("Customer validation failed")
    }
    
    // Validate nested address
    addressErrs := validation.ValidateStruct(customer.Address)
    if addressErrs.HasErrors() {
        fmt.Println("Address validation failed:")
        for _, err := range addressErrs {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Example 6: Custom Error Messages

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

func main() {
    data := map[string]any{
        "username": "abc",
        "age":      15,
    }
    
    rules := map[string][]validation.Rule{
        "username": {
            validation.Required("Please enter a username"),
            validation.StrMinLength(6, "Username must be at least 6 characters long"),
        },
        "age": {
            validation.NumMin(18, "You must be at least 18 years old to register"),
        },
    }
    
    errs := validation.ValidateMap(data, rules)
    
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Example 7: Conditional Validation

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

type PaymentRequest struct {
    Method      string `json:"method" validate:"required"`
    CardNumber  string `json:"card_number"`
    BankAccount string `json:"bank_account"`
}

func validatePayment(req PaymentRequest) validation.Errors {
    var errs validation.Errors
    
    // Validate method
    methodErr := validation.Validate("method", req.Method, "required")
    if methodErr != nil {
        errs = append(errs, *methodErr)
        return errs
    }
    
    // Conditional validation based on method
    if req.Method == "credit_card" {
        if req.CardNumber == "" {
            errs = append(errs, validation.Error{
                Field:   "card_number",
                Message: "Card number is required for credit card payment",
            })
        } else {
            cardErr := validation.Validate("card_number", req.CardNumber, "creditcard")
            if cardErr != nil {
                errs = append(errs, *cardErr)
            }
        }
    } else if req.Method == "bank_transfer" {
        if req.BankAccount == "" {
            errs = append(errs, validation.Error{
                Field:   "bank_account",
                Message: "Bank account is required for bank transfer",
            })
        }
    }
    
    return errs
}

func main() {
    req := PaymentRequest{
        Method:     "credit_card",
        CardNumber: "1234",  // Invalid
    }
    
    errs := validatePayment(req)
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Example 8: Batch Validation

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

type User struct {
    Name  string `json:"name" validate:"required,strminlen=3"`
    Email string `json:"email" validate:"required,email"`
}

func main() {
    users := []User{
        {Name: "John Doe", Email: "john@example.com"},
        {Name: "AB", Email: "invalid"},  // Invalid
        {Name: "Jane Smith", Email: "jane@example.com"},
        {Name: "", Email: "test@example.com"},  // Invalid
    }
    
    for i, user := range users {
        errs := validation.ValidateStruct(user)
        if errs.HasErrors() {
            fmt.Printf("User %d validation failed:\n", i+1)
            for _, err := range errs {
                fmt.Printf("  - %s: %s\n", err.Field, err.Message)
            }
        } else {
            fmt.Printf("User %d: Valid ✓\n", i+1)
        }
    }
}
```

---

## Custom Rules

You can create your own validation rules using the `Custom` function.

### Example: Custom Rule

```go
package main

import (
    "fmt"
    "strings"
    "github.com/fatkulnurk/foundation/validation"
)

// Custom rule: Check if string starts with a specific prefix
func startsWithPrefix(prefix string) validation.Rule {
    return validation.Custom(func(field string, value any) *validation.Error {
        s, ok := value.(string)
        if !ok {
            return &validation.Error{
                Field:   field,
                Message: "must be a string",
            }
        }
        
        if !strings.HasPrefix(s, prefix) {
            return &validation.Error{
                Field:   field,
                Message: fmt.Sprintf("must start with '%s'", prefix),
            }
        }
        
        return nil
    })
}

func main() {
    data := map[string]any{
        "product_code": "ABC123",
        "order_id":     "XYZ789",
    }
    
    rules := map[string][]validation.Rule{
        "product_code": {
            validation.Required(""),
            startsWithPrefix("PROD-"),  // Custom rule
        },
        "order_id": {
            validation.Required(""),
            startsWithPrefix("ORD-"),  // Custom rule
        },
    }
    
    errs := validation.ValidateMap(data, rules)
    
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Example: Complex Custom Rule

```go
package main

import (
    "fmt"
    "github.com/fatkulnurk/foundation/validation"
)

// Custom rule: Check if value is in a list
func oneOf(values []string) validation.Rule {
    return validation.Custom(func(field string, value any) *validation.Error {
        s, ok := value.(string)
        if !ok {
            return &validation.Error{
                Field:   field,
                Message: "must be a string",
            }
        }
        
        for _, v := range values {
            if s == v {
                return nil
            }
        }
        
        return &validation.Error{
            Field:   field,
            Message: fmt.Sprintf("must be one of: %v", values),
        }
    })
}

func main() {
    data := map[string]any{
        "status": "pending",
        "role":   "admin",
    }
    
    rules := map[string][]validation.Rule{
        "status": {
            validation.Required(""),
            oneOf([]string{"pending", "approved", "rejected"}),
        },
        "role": {
            validation.Required(""),
            oneOf([]string{"user", "admin", "moderator"}),
        },
    }
    
    errs := validation.ValidateMap(data, rules)
    
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    } else {
        fmt.Println("All validations passed!")
    }
}
```

---

## Best Practices

### 1. Use Struct Tags for Static Validation

```go
// Good - Clear and declarative
type User struct {
    Name  string `json:"name" validate:"required,strminlen=3"`
    Email string `json:"email" validate:"required,email"`
}

errs := validation.ValidateStruct(user)
```

### 2. Use Map Validation for Dynamic Data

```go
// Good - For form data, API requests, etc.
data := map[string]any{
    "name":  formData.Get("name"),
    "email": formData.Get("email"),
}

rules := map[string][]validation.Rule{
    "name":  {validation.Required(""), validation.StrMinLength(3, "")},
    "email": {validation.Required(""), validation.Email("")},
}

errs := validation.ValidateMap(data, rules)
```

### 3. Provide Custom Error Messages

```go
// Good - User-friendly messages
rules := map[string][]validation.Rule{
    "age": {
        validation.NumMin(18, "You must be at least 18 years old"),
    },
}

// Bad - Generic messages
rules := map[string][]validation.Rule{
    "age": {
        validation.NumMin(18, ""),  // Uses default message
    },
}
```

### 4. Validate Early

```go
// Good - Validate before processing
func CreateUser(req CreateUserRequest) error {
    errs := validation.ValidateStruct(req)
    if errs.HasErrors() {
        return fmt.Errorf("validation failed: %s", errs.Error())
    }
    
    // Process valid data
    return saveUser(req)
}
```

### 5. Group Related Validations

```go
// Good - Logical grouping
type UserProfile struct {
    // Personal info
    FirstName string `json:"first_name" validate:"required,strminlen=2"`
    LastName  string `json:"last_name" validate:"required,strminlen=2"`
    
    // Contact info
    Email string `json:"email" validate:"required,email"`
    Phone string `json:"phone" validate:"phone"`
    
    // Account info
    Username string `json:"username" validate:"required,username"`
    Password string `json:"password" validate:"required,password"`
}
```

### 6. Handle Errors Gracefully

```go
// Good - Detailed error handling
errs := validation.ValidateStruct(user)
if errs.HasErrors() {
    // Return structured errors for API
    return &APIError{
        Code:    "VALIDATION_ERROR",
        Message: "Validation failed",
        Fields:  errs,
    }
}
```

### 7. Use Constants for Validation Rules

```go
// Good - Reusable constants
const (
    MinUsernameLength = 6
    MaxUsernameLength = 20
    MinPasswordLength = 8
    MinAge            = 18
)

type User struct {
    Username string `json:"username" validate:"required,strminlen=6,strmaxlen=20"`
    Password string `json:"password" validate:"required,strminlen=8"`
    Age      int    `json:"age" validate:"nummin=18"`
}
```

### 8. Test Your Validations

```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {
            name:    "valid user",
            user:    User{Name: "John", Email: "john@example.com", Age: 25},
            wantErr: false,
        },
        {
            name:    "invalid email",
            user:    User{Name: "John", Email: "invalid", Age: 25},
            wantErr: true,
        },
        {
            name:    "too young",
            user:    User{Name: "John", Email: "john@example.com", Age: 15},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            errs := validation.ValidateStruct(tt.user)
            if (errs.HasErrors()) != tt.wantErr {
                t.Errorf("ValidateStruct() error = %v, wantErr %v", errs, tt.wantErr)
            }
        })
    }
}
```

---

## Database Constants

The package includes constants for common database field limits:

```go
const (
    DBVarcharMaxLength = 255      // VARCHAR max length
    DBTextMaxLength    = 65535    // TEXT max length
    DBTinyintMaxValue  = 127      // TINYINT max value
    DBSmallintMaxValue = 32767    // SMALLINT max value
    DBIntMaxValue      = 2147483647  // INT max value
    DBBigintMaxValue   = 9223372036854775807  // BIGINT max value
    DBFloatMaxDigits   = 7        // FLOAT precision
    DBDoubleMaxDigits  = 15       // DOUBLE precision
    DBDecimalMaxDigits = 65       // DECIMAL max digits
)
```

**Usage:**
```go
type User struct {
    Name string `json:"name" validate:"required,strmaxlen=255"`  // VARCHAR(255)
    Bio  string `json:"bio" validate:"strmaxlen=65535"`  // TEXT
}
```

---

## License

MIT

---

## Summary

The validation package provides a **simple, powerful way to validate data** in Go:
- Struct tag validation for static schemas
- Map validation for dynamic data
- 20+ built-in validation rules
- Custom rule support
- Clear error messages
- Type-safe and easy to use

**Key Features:**
- Required, min/max length, min/max value
- Email, phone, password, username validation
- URL, IP, UUID, JSON validation
- Credit card, postal code, hex color validation
- Custom rules for your specific needs

Now you can easily validate all your data in Go applications! 🚀
