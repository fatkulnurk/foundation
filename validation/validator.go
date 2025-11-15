package validation

import (
	"reflect"
	"strings"
)

func Validate(field string, value any, rule string) *Error {
	rules := parseTagToRules(rule)
	for _, rule := range rules {
		if err := rule(field, value); err != nil {
			return err
		}
	}

	return nil
}

// =========================
//  VALIDASI BERBASIS MAP
// =========================

// ValidateMap menjalankan rules terhadap data map.
//
// data: map[fieldName]value
// rules: map[fieldName][]Rule
func ValidateMap(data map[string]any, rules map[string][]Rule) Errors {
	var errs Errors

	for field, fieldRules := range rules {
		value, ok := data[field]
		if !ok {
			value = nil
		}

		for _, rule := range fieldRules {
			if err := rule(field, value); err != nil {
				// Selalu simpan field & message
				errs = append(errs, *err)
			}
		}
	}

	return errs
}

// =========================
//  VALIDASI BERBASIS STRUCT TAG
// =========================
//
// Tag yang didukung (di tag `validate`):
//   - validateRequired
//   - min=<angka>      -> untuk angka (pakai Rule validateNumMin)
//   - max=<angka>      -> untuk angka (pakai Rule validateNumMax)
//   - minlen=<angka>   -> untuk string (pakai Rule validateStrMinLength)
//   - maxlen=<angka>   -> untuk string (pakai Rule validateStrMaxLength)
//   - email            -> cek sederhana format email (harus mengandung "@")
//
// Penentuan nama field di error:
//   - pakai tag `json:"name"` kalau ada (nama sebelum koma)
//   - kalau tidak ada json tag, pakai nama field struct-nya

// ValidateStruct membaca tag `validate` pada struct dan mengembalikan Errors.
func ValidateStruct(s any) Errors {
	var errs Errors

	v := reflect.ValueOf(s)
	if !v.IsValid() {
		return errs
	}

	// Kalau pointer ke struct, dereference dulu
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return errs
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		// Tidak panic supaya aman, cukup return kosong
		return errs
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)

		// Skip field yang tidak bisa diakses (unexported)
		if !fieldValue.CanInterface() {
			continue
		}

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		// Tentukan nama field untuk error (json tag > nama struct)
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "" {
			fieldName = fieldType.Name
		} else {
			// json:"name,omitempty" -> ambil sebelum koma
			if idx := strings.Index(fieldName, ","); idx > 0 {
				fieldName = fieldName[:idx]
			}
		}

		rules := parseTagToRules(tag)

		val := fieldValue.Interface()
		for _, rule := range rules {
			if err := rule(fieldName, val); err != nil {
				errs = append(errs, *err)
			}
		}
	}

	return errs
}
