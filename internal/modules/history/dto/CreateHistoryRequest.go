package dtoHistory

import (
	"github.com/go-playground/validator/v10"
)

// Singleton validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type CreateHistoryRequest struct {
	Text   string  `json:"text" validate:"required,min=1"`
	Voice  string  `json:"voice" validate:"required"`
	Rate   float64 `json:"rate" validate:"min=0.1,max=5"`
	Pitch  float64 `json:"pitch" validate:"min=0,max=2"`
	Volume float64 `json:"volume" validate:"min=0,max=1"`
}

func (r *CreateHistoryRequest) Validate() error {
	return validate.Struct(r)
}

func ValidateCreateHistoryRequest(req *CreateHistoryRequest) error {
	return validate.Struct(req)
}
