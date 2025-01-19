package data

import "github.com/google/uuid"

type Product struct {
	ID         string
	Name       string
	ProviderID string
}

func NewProduct() Product {
	return Product{
		ID:         uuid.New().String(),
		Name:       "???",
		ProviderID: "",
	}
}
