package main

import (
	"fmt"
	"os"

	"github.com/fatkulnurk/foundation/validation"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "struct":
			runStructValidationExample()
		case "single":
			runSingleFieldExample()
		case "custom":
			runCustomRuleExample()
		case "api":
			runAPIExample()
		case "all":
			runAllRulesExample()
		default:
			fmt.Println("Usage: go run main.go [struct|single|custom|api|all]")
		}
	} else {
		runStructValidationExample()
	}
}

// runStructValidationExample demonstrates struct tag validation
func runStructValidationExample() {
	fmt.Println("=== Struct Validation Example ===")
	fmt.Println()

	type User struct {
		Name     string `json:"name" validate:"required,strminlen=3,strmaxlen=50"`
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,username"`
		Password string `json:"password" validate:"required,password"`
		Phone    string `json:"phone" validate:"required,phone"`
		Age      int    `json:"age" validate:"required,nummin=18,nummax=100"`
		Website  string `json:"website" validate:"url"`
	}

	// Example 1: Valid user
	fmt.Println("1. Validating valid user...")
	validUser := User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe123",
		Password: "SecureP@ss123",
		Phone:    "+1234567890",
		Age:      25,
		Website:  "https://example.com",
	}

	errs := validation.ValidateStruct(validUser)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ All validations passed!")
	}
	fmt.Println()

	// Example 2: Invalid user
	fmt.Println("2. Validating invalid user...")
	invalidUser := User{
		Name:     "Jo",            // Too short
		Email:    "invalid-email", // Invalid format
		Username: "abc",           // Too short
		Password: "weak",          // Too weak
		Phone:    "123",           // Invalid
		Age:      15,              // Too young
		Website:  "not-a-url",     // Invalid URL
	}

	errs = validation.ValidateStruct(invalidUser)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	}
	fmt.Println()

	// Example 3: Partially filled user
	fmt.Println("3. Validating partially filled user...")
	partialUser := User{
		Name:  "Jane Smith",
		Email: "jane@example.com",
		// Missing required fields
	}

	errs = validation.ValidateStruct(partialUser)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	}
	fmt.Println()

	fmt.Println("✅ Struct validation example completed!")
}

// runSingleFieldExample demonstrates single field validation
func runSingleFieldExample() {
	fmt.Println("=== Single Field Validation Example ===")
	fmt.Println()

	// Example 1: Validate email
	fmt.Println("1. Validating email addresses...")
	emails := []string{
		"valid@example.com",
		"invalid",
		"test@test.com",
		"no-at-sign",
	}

	for _, email := range emails {
		err := validation.Validate("email", email, "required,email")
		if err != nil {
			fmt.Printf("  ❌ '%s': %s\n", email, err.Message)
		} else {
			fmt.Printf("  ✅ '%s': Valid\n", email)
		}
	}
	fmt.Println()

	// Example 2: Validate age
	fmt.Println("2. Validating ages...")
	ages := []int{25, 15, 100, 150}

	for _, age := range ages {
		err := validation.Validate("age", age, "nummin=18,nummax=100")
		if err != nil {
			fmt.Printf("  ❌ Age %d: %s\n", age, err.Message)
		} else {
			fmt.Printf("  ✅ Age %d: Valid\n", age)
		}
	}
	fmt.Println()

	// Example 3: Validate username
	fmt.Println("3. Validating usernames...")
	usernames := []string{
		"johndoe123",
		"abc",
		"valid_user",
		"invalid-user",
	}

	for _, username := range usernames {
		err := validation.Validate("username", username, "username")
		if err != nil {
			fmt.Printf("  ❌ '%s': %s\n", username, err.Message)
		} else {
			fmt.Printf("  ✅ '%s': Valid\n", username)
		}
	}
	fmt.Println()

	fmt.Println("✅ Single field validation example completed!")
}

// runCustomRuleExample demonstrates custom validation rules
func runCustomRuleExample() {
	fmt.Println("=== Custom Rule Example ===")
	fmt.Println()

	// Custom rule: Check if string starts with prefix
	startsWithPrefix := func(prefix string) validation.Rule {
		return validation.Custom(func(field string, value any) *validation.Error {
			s, ok := value.(string)
			if !ok {
				return &validation.Error{
					Field:   field,
					Message: "must be a string",
				}
			}

			if len(s) < len(prefix) || s[:len(prefix)] != prefix {
				return &validation.Error{
					Field:   field,
					Message: fmt.Sprintf("must start with '%s'", prefix),
				}
			}

			return nil
		})
	}

	// Custom rule: Check if value is in list
	oneOf := func(values []string) validation.Rule {
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

	// Example 1: Valid data
	fmt.Println("1. Validating with custom rules (valid)...")
	validData := map[string]any{
		"order_id": "ORD-12345",
		"status":   "pending",
		"priority": "high",
	}

	customRules := map[string][]validation.Rule{
		"order_id": {
			startsWithPrefix("ORD-"),
		},
		"status": {
			oneOf([]string{"pending", "processing", "completed", "cancelled"}),
		},
		"priority": {
			oneOf([]string{"low", "medium", "high", "urgent"}),
		},
	}

	errs := validation.ValidateMap(validData, customRules)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ All custom validations passed!")
	}
	fmt.Println()

	// Example 2: Invalid data
	fmt.Println("2. Validating with custom rules (invalid)...")
	invalidData := map[string]any{
		"order_id": "12345",      // Missing prefix
		"status":   "unknown",    // Not in list
		"priority": "super-high", // Not in list
	}

	errs = validation.ValidateMap(invalidData, customRules)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	}
	fmt.Println()

	fmt.Println("✅ Custom rule example completed!")
}

// runAPIExample demonstrates API request validation
func runAPIExample() {
	fmt.Println("=== API Request Validation Example ===")
	fmt.Println()

	type CreateUserRequest struct {
		Name     string `json:"name" validate:"required,strminlen=3,strmaxlen=50"`
		Email    string `json:"email" validate:"required,email"`
		Username string `json:"username" validate:"required,username"`
		Password string `json:"password" validate:"required,password"`
		Phone    string `json:"phone" validate:"phone"`
		Age      int    `json:"age" validate:"nummin=18,nummax=100"`
	}

	type UpdateProfileRequest struct {
		Name    string `json:"name" validate:"strminlen=3,strmaxlen=50"`
		Bio     string `json:"bio" validate:"strmaxlen=500"`
		Website string `json:"website" validate:"url"`
	}

	// Example 1: Create user request (valid)
	fmt.Println("1. Validating create user request (valid)...")
	createReq := CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe123",
		Password: "SecureP@ss123",
		Phone:    "+1234567890",
		Age:      25,
	}

	errs := validation.ValidateStruct(createReq)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ Create user request is valid!")
	}
	fmt.Println()

	// Example 2: Create user request (invalid)
	fmt.Println("2. Validating create user request (invalid)...")
	invalidCreateReq := CreateUserRequest{
		Name:     "Jo",
		Email:    "invalid",
		Username: "abc",
		Password: "weak",
		Age:      15,
	}

	errs = validation.ValidateStruct(invalidCreateReq)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	}
	fmt.Println()

	// Example 3: Update profile request
	fmt.Println("3. Validating update profile request...")
	updateReq := UpdateProfileRequest{
		Name:    "Jane Smith",
		Bio:     "Software developer passionate about Go",
		Website: "https://janesmith.dev",
	}

	errs = validation.ValidateStruct(updateReq)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ Update profile request is valid!")
	}
	fmt.Println()

	fmt.Println("✅ API validation example completed!")
}

// runAllRulesExample demonstrates all available validation rules
func runAllRulesExample() {
	fmt.Println("=== All Validation Rules Example ===")
	fmt.Println()

	type AllRulesExample struct {
		// String rules
		RequiredField string `json:"required_field" validate:"required"`
		MinLenField   string `json:"min_len_field" validate:"strminlen=5"`
		MaxLenField   string `json:"max_len_field" validate:"strmaxlen=10"`
		EmailField    string `json:"email_field" validate:"email"`
		UsernameField string `json:"username_field" validate:"username"`
		PasswordField string `json:"password_field" validate:"password"`
		PhoneField    string `json:"phone_field" validate:"phone"`
		URLField      string `json:"url_field" validate:"url"`
		AlphaNumField string `json:"alphanum_field" validate:"alphanumeric"`

		// Number rules
		MinNumField int     `json:"min_num_field" validate:"nummin=10"`
		MaxNumField int     `json:"max_num_field" validate:"nummax=100"`
		AgeField    int     `json:"age_field" validate:"nummin=18,nummax=100"`
		PriceField  float64 `json:"price_field" validate:"nummin=0"`

		// Format rules
		DateField       string `json:"date_field" validate:"date"`
		UUIDField       string `json:"uuid_field" validate:"uuid"`
		JSONField       string `json:"json_field" validate:"json"`
		HexColorField   string `json:"hex_color_field" validate:"hexcolor"`
		CreditCardField string `json:"credit_card_field" validate:"creditcard"`
		PostalCodeField string `json:"postal_code_field" validate:"postalcode"`
		Base64Field     string `json:"base64_field" validate:"base64"`

		// Network rules
		IPField   string `json:"ip_field" validate:"ip"`
		IPv4Field string `json:"ipv4_field" validate:"ipv4"`
		IPv6Field string `json:"ipv6_field" validate:"ipv6"`
	}

	// Valid example
	fmt.Println("1. Testing all rules with valid data...")
	validExample := AllRulesExample{
		RequiredField:   "present",
		MinLenField:     "12345",
		MaxLenField:     "short",
		EmailField:      "test@example.com",
		UsernameField:   "johndoe123",
		PasswordField:   "SecureP@ss123",
		PhoneField:      "+1234567890",
		URLField:        "https://example.com",
		AlphaNumField:   "ABC123",
		MinNumField:     50,
		MaxNumField:     75,
		AgeField:        25,
		PriceField:      99.99,
		DateField:       "2024-11-26",
		UUIDField:       "550e8400-e29b-41d4-a716-446655440000",
		JSONField:       `{"key":"value"}`,
		HexColorField:   "#FF5733",
		CreditCardField: "4532015112830366",
		PostalCodeField: "12345",
		Base64Field:     "SGVsbG8gV29ybGQ=",
		IPField:         "192.168.1.1",
		IPv4Field:       "192.168.1.1",
		IPv6Field:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
	}

	errs := validation.ValidateStruct(validExample)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed:")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ All rules passed with valid data!")
	}
	fmt.Println()

	// Invalid example
	fmt.Println("2. Testing all rules with invalid data...")
	invalidExample := AllRulesExample{
		RequiredField:   "",
		MinLenField:     "123",
		MaxLenField:     "this is too long",
		EmailField:      "invalid",
		UsernameField:   "abc",
		PasswordField:   "weak",
		PhoneField:      "123",
		URLField:        "not-a-url",
		AlphaNumField:   "ABC-123",
		MinNumField:     5,
		MaxNumField:     150,
		AgeField:        15,
		PriceField:      -10,
		DateField:       "2024/11/26",
		UUIDField:       "invalid-uuid",
		JSONField:       "not json",
		HexColorField:   "FF5733",
		CreditCardField: "1234",
		PostalCodeField: "123",
		Base64Field:     "not base64!",
		IPField:         "999.999.999.999",
		IPv4Field:       "invalid",
		IPv6Field:       "invalid",
	}

	errs = validation.ValidateStruct(invalidExample)
	if errs.HasErrors() {
		fmt.Println("❌ Validation failed (as expected):")
		for _, err := range errs {
			fmt.Printf("  - %s: %s\n", err.Field, err.Message)
		}
	}
	fmt.Println()

	fmt.Println("✅ All rules example completed!")
}
