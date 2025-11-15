package validation

import (
	"net"
	"strconv"
	"strings"
)

// parseTagToRules mengubah string tag validate jadi slice Rule.
//
// Contoh:
//
//	"validateRequired,minlen=3,maxlen=10,email"
//	"min=18,max=60"
func parseTagToRules(tag string) []Rule {
	parts := strings.Split(tag, ",")
	var rules []Rule

	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}

		switch {
		case p == RuleRequired:
			rules = append(rules, validateRequired(""))

		case strings.HasPrefix(p, RuleStrMinLength):
			nStr := strings.TrimPrefix(p, RuleStrMinLength)
			n, err := strconv.Atoi(nStr)
			if err == nil {
				rules = append(rules, validateStrMinLength(n, ""))
			}

		case strings.HasPrefix(p, RuleStrMaxLength):
			nStr := strings.TrimPrefix(p, RuleStrMaxLength)
			n, err := strconv.Atoi(nStr)
			if err == nil {
				rules = append(rules, validateStrMaxLength(n, ""))
			}

		case strings.HasPrefix(p, RuleNumMin):
			nStr := strings.TrimPrefix(p, RuleNumMin)
			n, err := strconv.ParseFloat(nStr, 64)
			if err == nil {
				rules = append(rules, validateNumMin(n, ""))
			}

		case strings.HasPrefix(p, RuleNumMax):
			nStr := strings.TrimPrefix(p, RuleNumMax)
			n, err := strconv.ParseFloat(nStr, 64)
			if err == nil {
				rules = append(rules, validateNumMax(n, ""))
			}

		case p == RuleEmail:
			rules = append(rules, Custom(func(field string, value any) *Error {
				s, ok := value.(string)
				if !ok {
					return &Error{
						Field:   field,
						Message: "harus berupa string (email)",
					}
				}
				s = strings.TrimSpace(s)
				if s == "" {
					// Biar rule validateRequired yang handle kalau dipakai
					return nil
				}
				if !strings.Contains(s, "@") {
					return &Error{
						Field:   field,
						Message: "format email tidak valid",
					}
				}
				return nil
			}))
		case p == RuleIP:
			rules = append(rules, Custom(func(field string, value any) *Error {
				s, ok := value.(string)
				if !ok {
					return &Error{
						Field:   field,
						Message: ErrorMessageString,
					}
				}
				s = strings.TrimSpace(s)
				if s == "" {
					// Biar rule validateRequired yang handle kalau dipakai
					return nil
				}
				ip := net.ParseIP(s)
				if ip == nil {
					return &Error{
						Field:   field,
						Message: ErrorMessageInvalidIP,
					}
				}
				return nil
			}))
		case p == RuleIPv4:
			rules = append(rules, Custom(func(field string, value any) *Error {
				s, ok := value.(string)
				if !ok {
					return &Error{
						Field:   field,
						Message: ErrorMessageString,
					}
				}
				s = strings.TrimSpace(s)
				if s == "" {
					// Biar rule validateRequired yang handle kalau dipakai
					return nil
				}
				ip := net.ParseIP(s)
				if ip == nil || ip.To4() == nil {
					return &Error{
						Field:   field,
						Message: ErrorMessageInvalidIPv4,
					}
				}
				return nil
			}))
		case p == RuleIPv6:
			rules = append(rules, Custom(func(field string, value any) *Error {
				s, ok := value.(string)
				if !ok {
					return &Error{
						Field:   field,
						Message: ErrorMessageString,
					}
				}
				s = strings.TrimSpace(s)
				if s == "" {
					// Biar rule validateRequired yang handle kalau dipakai
					return nil
				}
				ip := net.ParseIP(s)
				if ip == nil || ip.To16() == nil {
					return &Error{
						Field:   field,
						Message: ErrorMessageInvalidIPv6,
					}
				}
				return nil
			}))
		case p == RuleUsername:
			rules = append(rules, validateUsername(""))
		case p == RulePhone:
			rules = append(rules, validatePhone(""))
		case p == RulePassword:
			rules = append(rules, validatePassword(""))
		case p == RuleURL:
			rules = append(rules, validateURLFormat(""))
		case p == RuleDate:
			rules = append(rules, validateDate(""))
		case p == RuleAlphaNumeric:
			rules = append(rules, validateAlphaNumeric(""))
		case p == RuleUUID:
			rules = append(rules, validateUuid(""))
		case p == RuleJSON:
			rules = append(rules, validateJson(""))
		case p == RuleHexColor:
			rules = append(rules, validateHexColor(""))
		case p == RuleCreditCard:
			rules = append(rules, validateCreditCard(""))
		case p == RulePostalCode:
			rules = append(rules, validatePostalCode(""))
		case p == RuleBase64:
			rules = append(rules, validateBase64(""))
		default:
			// Tag tidak dikenal
			panic("unknown tag: " + p)
		}
	}

	return rules
}

// helper: konversi berbagai tipe angka ke float64.
func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}
