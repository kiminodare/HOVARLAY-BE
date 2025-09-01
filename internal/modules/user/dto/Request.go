package dtoUser

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

func (r *Request) Validate() error {
	return validate.Struct(r)
}

func ValidateRequest(req *Request) error {
	return validate.Struct(req)
}
