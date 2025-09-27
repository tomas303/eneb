package data

import (
	"github.com/google/uuid"
)

type Price struct {
	ID          string
	EnergyKind  int
	PriceType   int
	Value       int
	Provider_ID string
	Name        string
}

func NewPrice() Price {
	return Price{
		ID:          uuid.NewString(),
		EnergyKind:  0,
		PriceType:   0,
		Value:       0,
		Provider_ID: "",
		Name:        "???",
	}
}
