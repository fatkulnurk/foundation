package main

import (
	"fmt"

	"github.com/fatkulnurk/gostarter/pkg/validation"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	Name      string `json:"name"     validate:"required,strminlen=3,strmaxlen=50"`
	Email     string `json:"email"    validate:"required,email"`
	Age       int    `json:"age"      validate:"nummin=18,nummax=60"`
	Username  string
	Bio       string `json:"bio"      validate:"strmaxlen=200"`
}

func main() {
	req := RegisterRequest{
		Name:     "Fa",
		Email:    "salahformat",
		Age:      15,
		Username: "username_panjang_banget",
		Bio:      "ini bio singkat",
	}

	errs := validation.ValidateStruct(req)

	if errs.HasErrors() {
		fmt.Println("Validasi gagal:")
		for _, e := range errs {
			fmt.Printf("- %s: %s\n", e.Field, e.Message)
		}
	} else {
		fmt.Println("Validasi berhasil tanpa error")
	}
}
