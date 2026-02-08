package data

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          string
	EnergyKind  int
	PriceType   int
	Provider_ID string
	Name        string
}

func NewProduct() Product {
	return Product{
		ID:          uuid.NewString(),
		EnergyKind:  0,
		PriceType:   0,
		Provider_ID: "",
		Name:        "???",
	}
}
