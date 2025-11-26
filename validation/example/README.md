# Validation Package Examples

This directory contains working examples demonstrating how to use the validation package.

## Prerequisites

1. **Go 1.25 or higher**
2. **Dependencies:** `github.com/google/uuid`

## Running the Examples

### 1. Struct Validation Example

Demonstrates validation using struct tags.

```bash
cd pkg/validation/example
go run main.go struct
```

**What it does:**
- Validates complete user data with all fields
- Shows validation of invalid data
- Demonstrates partial validation (missing fields)

**Output:**
```
=== Struct Validation Example ===

1. Validating valid user...
✅ All validations passed!

2. Validating invalid user...
❌ Validation failed:
  - name: must be at least 3 characters
  - email: format email tidak valid
  - username: must be between 6.00 and 16.00
  - password: must be between 8.00 and 16.00
  - phone: must be between 10.00 and 15.00
  - age: must be at least 18.00
  - url_field: is not a valid URL

3. Validating partially filled user...
❌ Validation failed:
  - username: can't be empty
  - password: can't be empty
  - phone: can't be empty
  - age: can't be empty
```

### 2. Single Field Validation Example

Demonstrates validating individual fields.

```bash
go run main.go single
```

**What it does:**
- Validates email addresses
- Validates ages with min/max rules
- Validates usernames

### 3. Custom Rule Example

Demonstrates creating and using custom validation rules.

```bash
go run main.go custom
```

**What it does:**
- Creates custom "starts with prefix" rule
- Creates custom "one of" rule
- Validates data with custom rules

### 4. API Request Validation Example

Demonstrates validation for API requests.

```bash
go run main.go api
```

**What it does:**
- Validates create user requests
- Validates update profile requests
- Shows how to use validation in API handlers

### 5. All Rules Example

Demonstrates all available validation rules.

```bash
go run main.go all
```

**What it does:**
- Tests all 20+ built-in validation rules
- Shows valid and invalid examples for each rule
- Comprehensive demonstration of all features

## Example Output

### Struct Validation (Valid Data)

```
=== Struct Validation Example ===

1. Validating valid user...
✅ All validations passed!
```

### Struct Validation (Invalid Data)

```
2. Validating invalid user...
❌ Validation failed:
  - name: must be at least 3 characters
  - email: format email tidak valid
  - username: must be between 6.00 and 16.00
  - password: must be between 8.00 and 16.00
  - phone: must be between 10.00 and 15.00
  - age: must be at least 18.00
  - website: is not a valid URL
```

### Single Field Validation

```
=== Single Field Validation Example ===

1. Validating email addresses...
  ✅ 'valid@example.com': Valid
  ❌ 'invalid': format email tidak valid
  ✅ 'test@test.com': Valid
  ❌ 'no-at-sign': format email tidak valid

2. Validating ages...
  ✅ Age 25: Valid
  ❌ Age 15: must be at least 18.00
  ✅ Age 100: Valid
  ❌ Age 150: must be at most 100.00
```

### Custom Rules

```
=== Custom Rule Example ===

1. Validating with custom rules (valid)...
✅ All custom validations passed!

2. Validating with custom rules (invalid)...
❌ Validation failed:
  - order_id: must start with 'ORD-'
  - status: must be one of: [pending processing completed cancelled]
  - priority: must be one of: [low medium high urgent]
```

## Code Examples

### Example 1: Basic Struct Validation

```go
type User struct {
    Name  string `json:"name" validate:"required,strminlen=3"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"nummin=18,nummax=100"`
}

user := User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   25,
}

errs := validation.ValidateStruct(user)
if errs.HasErrors() {
    for _, err := range errs {
        fmt.Printf("%s: %s\n", err.Field, err.Message)
    }
}
```

### Example 2: Single Field Validation

```go
// Validate email
err := validation.Validate("email", "test@example.com", "required,email")
if err != nil {
    fmt.Println(err.Message)
}

// Validate age
err = validation.Validate("age", 25, "nummin=18,nummax=100")
if err != nil {
    fmt.Println(err.Message)
}
```

### Example 3: Custom Rule

```go
// Create custom rule
startsWithPrefix := func(prefix string) validation.Rule {
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

// Use custom rule
data := map[string]any{
    "order_id": "ORD-12345",
}

rules := map[string][]validation.Rule{
    "order_id": {startsWithPrefix("ORD-")},
}

errs := validation.ValidateMap(data, rules)
```

## Available Validation Rules

### String Rules
- `required` - Cannot be empty
- `strminlen=N` - Minimum length
- `strmaxlen=N` - Maximum length
- `email` - Valid email format
- `username` - Valid username (6-16 chars, letters, numbers, underscores)
- `password` - Strong password (8-16 chars, uppercase, lowercase, number, special char)
- `phone` - Valid phone number (10-15 digits)
- `url` - Valid URL
- `alphanumeric` - Only letters and numbers

### Number Rules
- `nummin=N` - Minimum value
- `nummax=N` - Maximum value

### Format Rules
- `date` - Valid date (YYYY-MM-DD)
- `uuid` - Valid UUID
- `json` - Valid JSON string
- `hexcolor` - Valid hex color (#RRGGBB)
- `creditcard` - Valid credit card (Luhn algorithm)
- `postalcode` - Valid postal code (5 digits)
- `base64` - Valid base64 string

### Network Rules
- `ip` - Valid IP address
- `ipv4` - Valid IPv4 address
- `ipv6` - Valid IPv6 address

## Modifying the Examples

### Add Your Own Validation

```go
type MyStruct struct {
    CustomField string `json:"custom_field" validate:"required,strminlen=5"`
}

func runMyExample() {
    data := MyStruct{
        CustomField: "test",
    }
    
    errs := validation.ValidateStruct(data)
    if errs.HasErrors() {
        for _, err := range errs {
            fmt.Printf("%s: %s\n", err.Field, err.Message)
        }
    }
}
```

Then add to main():
```go
case "myexample":
    runMyExample()
```

## Common Patterns

### API Request Validation

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    errs := validation.ValidateStruct(req)
    if errs.HasErrors() {
        json.NewEncoder(w).Encode(map[string]any{
            "error": "Validation failed",
            "fields": errs,
        })
        return
    }
    
    // Process valid request
}
```

### Form Validation

```go
func validateForm(formData map[string]string) validation.Errors {
    data := map[string]any{
        "name":  formData["name"],
        "email": formData["email"],
        "age":   formData["age"],
    }
    
    rules := map[string][]validation.Rule{
        "name":  {/* rules */},
        "email": {/* rules */},
        "age":   {/* rules */},
    }
    
    return validation.ValidateMap(data, rules)
}
```

### Conditional Validation

```go
func validatePayment(req PaymentRequest) validation.Errors {
    var errs validation.Errors
    
    // Always validate method
    if err := validation.Validate("method", req.Method, "required"); err != nil {
        errs = append(errs, *err)
    }
    
    // Conditional validation
    if req.Method == "credit_card" {
        if err := validation.Validate("card_number", req.CardNumber, "creditcard"); err != nil {
            errs = append(errs, *err)
        }
    }
    
    return errs
}
```

## Troubleshooting

### Error: "unknown tag"

**Problem:** Using an unsupported validation tag

**Solution:** Check the list of available rules in the main README.md

### Error: "must be a string"

**Problem:** Applying string validation to non-string field

**Solution:** Ensure field type matches the validation rule

### Validation Not Working

**Problem:** Validation passes but shouldn't

**Solution:**
- Check if field is exported (starts with capital letter)
- Verify the `validate` tag syntax
- Ensure `json` tag is present for field name mapping

## Next Steps

After running these examples:

1. **Integrate into your application**
   - Copy patterns from examples
   - Adapt to your data structures

2. **Create custom rules**
   - Use `validation.Custom()` for specific needs
   - Combine with built-in rules

3. **Test thoroughly**
   - Write unit tests for validations
   - Test edge cases

## Learn More

See the main [README.md](../README.md) for:
- Complete API documentation
- All validation rules
- Best practices
- Advanced usage

---

Happy validating! 🚀
