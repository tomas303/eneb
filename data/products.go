package data

import "github.com/google/uuid"

type Product struct {
	ID          string
	Name        string
	Provider_ID string
	EnergyKind  int
}

func NewProduct() Product {
	return Product{
		ID:          uuid.New().String(),
		Name:        "???",
		Provider_ID: "",
		EnergyKind:  0,
	}
}
