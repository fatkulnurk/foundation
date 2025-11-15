package validation

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Rule adalah fungsi yang memvalidasi satu field.
type Rule func(field string, value any) *Error

// Custom : membuat rule kustom sendiri.
func Custom(f func(field string, value any) *Error) Rule {
	return func(field string, value any) *Error {
		return f(field, value)
	}
}

// =========================
//  RULE BUILDER UMUM
// =========================

// validateRequired : value tidak boleh kosong (nil atau string kosong / spasi).
// Jika message kosong (""), akan pakai default message.
func validateRequired(message string) Rule {
	return func(field string, value any) *Error {
		msg := message
		if value == nil {
			if msg == "" {
				msg = ErrorMessageRequired
			}
			return &Error{Field: field, Message: msg}
		}

		// Jika string, cek kosong / spasi
		if s, ok := value.(string); ok {
			if strings.TrimSpace(s) == "" {
				if msg == "" {
					msg = ErrorMessageRequired
				}
				return &Error{Field: field, Message: msg}
			}
		}

		return nil
	}
}

// validateStrMinLength : panjang minimal untuk string.
func validateStrMinLength(min int, message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		if len(s) < min {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageMinLength, min)
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

// validateStrMaxLength : panjang maksimal untuk string.
func validateStrMaxLength(max int, message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		if len(s) > max {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageMaxLength, max)
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

// validateNumMin : nilai minimal untuk angka (int, uint, float32, float64, dll).
func validateNumMin(min float64, message string) Rule {
	return func(field string, value any) *Error {
		if value == nil {
			// Biarkan validateRequired yang handle kalau dipakai
			return nil
		}

		num, ok := toFloat(value)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageNumber,
			}
		}

		if num < min {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageMin, min)
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

// validateNumMax : nilai maksimal untuk angka (int, uint, float32, float64, dll).
func validateNumMax(max float64, message string) Rule {
	return func(field string, value any) *Error {
		if value == nil {
			return nil
		}

		num, ok := toFloat(value)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageNumber,
			}
		}

		if num > max {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageMax, max)
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validateUsername(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		if len(s) < 6 || len(s) > 16 {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageNumberBetween, float64(6), float64(16))
			}
			return &Error{Field: field, Message: msg}
		}

		for _, ch := range s {
			if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_') {
				return &Error{
					Field:   field,
					Message: ErrorMessageInvalidUsername,
				}
			}
		}

		return nil
	}
}

func validatePhone(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		digitCount := 0
		for i, ch := range s {
			// hitung digit
			if ch >= '0' && ch <= '9' {
				digitCount++
				continue
			}

			// karakter non-digit yang diizinkan
			if ch == ' ' || ch == '-' || ch == '(' || ch == ')' {
				continue
			}
			if ch == '+' && i == 0 {
				continue
			}

			// selain itu dianggap tidak valid
			return &Error{
				Field:   field,
				Message: ErrorMessageInvalidPhone,
			}
		}

		if digitCount < 10 || digitCount > 15 {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageNumberBetween, float64(10), float64(15))
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validatePassword(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		if len(s) < 8 || len(s) > 16 {
			msg := message
			if msg == "" {
				msg = fmt.Sprintf(ErrorMessageNumberBetween, float64(8), float64(16))
			}
			return &Error{Field: field, Message: msg}
		}

		hasUpper := false
		hasLower := false
		hasDigit := false
		hasSpecial := false

		for _, ch := range s {
			switch {
			case ch >= 'A' && ch <= 'Z':
				hasUpper = true
			case ch >= 'a' && ch <= 'z':
				hasLower = true
			case ch >= '0' && ch <= '9':
				hasDigit = true
			default:
				hasSpecial = true
			}
		}

		if !(hasUpper && hasLower && hasDigit && hasSpecial) {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidPassword
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validateURLFormat(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		u, err := url.Parse(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidURL
			}
			return &Error{Field: field, Message: msg}
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidURL
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validateDate(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		_, err := time.Parse("2006-01-02", s)
		if err != nil {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidDate
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validateAlphaNumeric(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		for _, ch := range s {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
				msg := message
				if msg == "" {
					msg = ErrorMessageInvalidAlphaNumeric
				}
				return &Error{Field: field, Message: msg}
			}
		}

		return nil
	}
}

func validateUuid(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		_, err := uuid.Parse(s)
		if err != nil {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidUUID
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validateJson(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		if !json.Valid([]byte(s)) {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidJSON
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validateHexColor(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		if !strings.HasPrefix(s, "#") || len(s) != 7 {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidHexColor
			}
			return &Error{Field: field, Message: msg}
		}

		for _, ch := range s[1:] {
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				msg := message
				if msg == "" {
					msg = ErrorMessageInvalidHexColor
				}
				return &Error{Field: field, Message: msg}
			}
		}

		return nil
	}
}

func validateCreditCard(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		// normalisasi: ambil hanya digit
		digits := make([]int, 0, len(s))
		for _, ch := range s {
			if ch >= '0' && ch <= '9' {
				digits = append(digits, int(ch-'0'))
			} else if ch == ' ' || ch == '-' {
				// diabaikan (formatting saja)
				continue
			} else {
				// karakter lain tidak diizinkan
				msg := message
				if msg == "" {
					msg = ErrorMessageInvalidCreditCard
				}
				return &Error{Field: field, Message: msg}
			}
		}

		// panjang tipikal kartu kredit 13–19 digit
		if len(digits) < 13 || len(digits) > 19 {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidCreditCard
			}
			return &Error{Field: field, Message: msg}
		}

		// algoritma Luhn
		sum := 0
		double := false
		for i := len(digits) - 1; i >= 0; i-- {
			d := digits[i]
			if double {
				d *= 2
				if d > 9 {
					d -= 9
				}
			}
			sum += d
			double = !double
		}

		if sum%10 != 0 {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidCreditCard
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}

func validatePostalCode(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		// contoh: kode pos Indonesia = 5 digit
		if len(s) != 5 {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidPostalCode
			}
			return &Error{Field: field, Message: msg}
		}

		for _, ch := range s {
			if ch < '0' || ch > '9' {
				msg := message
				if msg == "" {
					msg = ErrorMessageInvalidPostalCode
				}
				return &Error{Field: field, Message: msg}
			}
		}

		return nil
	}
}

func validateBase64(message string) Rule {
	return func(field string, value any) *Error {
		s, ok := value.(string)
		if !ok {
			return &Error{
				Field:   field,
				Message: ErrorMessageString,
			}
		}

		s = strings.TrimSpace(s)
		if s == "" {
			// Biarkan validateRequired yang handle jika dipakai
			return nil
		}

		_, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			msg := message
			if msg == "" {
				msg = ErrorMessageInvalidBase64
			}
			return &Error{Field: field, Message: msg}
		}

		return nil
	}
}
