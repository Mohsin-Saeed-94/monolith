package pet

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Dto struct {
	Id     uuid.UUID  `json:"id" db:"id"`
	UserId *uuid.UUID `json:"userId" validate:"required" db:"user_id"`
	Breed  string     `json:"breed" db:"breed"`
	Weight float64    `json:"weight" db:"weight"`
}

func New() *Dto {
	return &Dto{
		Id:     uuid.New(),
		UserId: nil,
		Breed:  "",
		Weight: 0.00,
	}
}

func Read(r io.Reader) (*Dto, error) {
	d := Dto{}
	err := json.NewDecoder(r).Decode(&d)
	return &d, err
}

func (d *Dto) Validate() error {
	v := validator.New()
	return v.Struct(d)
}

func (d *Dto) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(d)
}
