package validation

import (
	"strings"
)

type Error struct {
	Field   string `json:"field" xml:"field" bson:"field"`
	Message string `json:"message" xml:"message" bson:"message"`
}

type Errors []Error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var b strings.Builder
	for i, e := range e {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Field)
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	return b.String()
}

func (e Errors) HasErrors() bool {
	return len(e) > 0
}

func (e Errors) ForField(field string) []Error {
	var result []Error
	for _, e := range e {
		if e.Field == field {
			result = append(result, e)
		}
	}
	return result
}
