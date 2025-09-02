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
	Text   string  `json:"text" validate:"required,min=1" custom:"sentence"`
	Voice  string  `json:"voice" validate:"required" custom:"word"`
	Rate   float64 `json:"rate" validate:"min=0.1,max=5" custom:"oneof: 0.5,1.0,1.5,2.0"`
	Pitch  float64 `json:"pitch" validate:"min=0,max=2"  custom:"oneof: 0.5,1.0,1.5,2.0"`
	Volume float64 `json:"volume" validate:"min=0,max=1"  custom:"oneof: 0.5,1.0"`
}

func (r *CreateHistoryRequest) Validate() error {
	return validate.Struct(r)
}
