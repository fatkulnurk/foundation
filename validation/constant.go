package validation

const (
	RuleRequired     = "validateRequired"
	RuleStrMinLength = "strminlen="
	RuleStrMaxLength = "strmaxlen="
	RuleNumMin       = "nummin="
	RuleNumMax       = "nummax="
	RuleEmail        = "email"
	RulePhone        = "validatePhone"
	RuleUsername     = "validateUsername"
	RulePassword     = "validatePassword"
	RuleURL          = "url"
	RuleDate         = "validateDate"
	RuleAlphaNumeric = "alphanumeric"
	RuleUUID         = "uuid"
	RuleJSON         = "json"
	RuleHexColor     = "hexcolor"
	RuleCreditCard   = "creditcard"
	RulePostalCode   = "postalcode"
	RuleBase64       = "base64"
	RuleIP           = "ip"
	RuleIPv4         = "ipv4"
	RuleIPv6         = "ipv6"
)

const (
	ErrorMessageRequired            = "can't be empty"
	ErrorMessageInvalidEmail        = "is not a valid email"
	ErrorMessageInvalidPhone        = "is not a valid validatePhone number"
	ErrorMessageInvalidPassword     = "must be at least 8 characters and include uppercase, lowercase, number, and special character"
	ErrorMessageInvalidUsername     = "must be 6-16 characters and only contain letters, numbers, and underscores"
	ErrorMessageInvalidURL          = "is not a valid URL"
	ErrorMessageInvalidIP           = "is not a valid IP address"
	ErrorMessageInvalidIPv4         = "is not a valid IPv4 address"
	ErrorMessageInvalidIPv6         = "is not a valid IPv6 address"
	ErrorMessageInvalidDate         = "is not a valid validateDate (YYYY-MM-DD)"
	ErrorMessageInvalidAlphaNumeric = "may only contain letters and numbers"
	ErrorMessageInvalidUUID         = "is not a valid UUID"
	ErrorMessageInvalidJSON         = "is not valid validateJson"
	ErrorMessageInvalidHexColor     = "is not a valid hex color code"
	ErrorMessageInvalidCreditCard   = "is not a valid credit card number"
	ErrorMessageInvalidPostalCode   = "is not a valid postal code"
	ErrorMessageInvalidBase64       = "is not a valid base64 string"
	ErrorMessageOneOf               = "must be one of %s"
	ErrorMessageLowercase           = "must be lowercase"
	ErrorMessageUppercase           = "must be uppercase"
	ErrorMessageString              = "must be a string"
	ErrorMessageNumber              = "must be a number"
	ErrorMessageFloat               = "must be a floating-point number"
	ErrorMessageInt                 = "must be an integer"
	ErrorMessageMin                 = "must be at least %.2f"
	ErrorMessageMax                 = "must be at most %.2f"
	ErrorMessageMinLength           = "must be at least %d characters"
	ErrorMessageMaxLength           = "must be at most %d characters"
	ErrorMessageGreaterThan         = "must be greater than %d"
	ErrorMessageLessThan            = "must be less than %d"
	ErrorMessageGreaterThanEqual    = "must be at least %d"
	ErrorMessageLessThanEqual       = "must be at most %d"
	ErrorMessageGreaterThanEqualStr = "must be at least %s"
	ErrorMessageLessThanEqualStr    = "must be at most %s"
	ErrorMessageNotEqual            = "must not be equal to %s"
	ErrorMessageNotEqualStr         = "must not be equal to %s"
	ErrorMessageNotIn               = "must not be in %s"
	ErrorMessageNumberBetween       = "must be between %.2f and %.2f"
)

const (
	DBVarcharMaxLength = 255
	DBTextMaxLength    = 65535
	DBTinyintMaxValue  = 127
	DBSmallintMaxValue = 32767
	DBIntMaxValue      = 2147483647
	DBBigintMaxValue   = 9223372036854775807
	DBFloatMaxDigits   = 7
	DBDoubleMaxDigits  = 15
	DBDecimalMaxDigits = 65
)
