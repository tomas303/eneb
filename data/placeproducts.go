package data

import (
	"time"

	"github.com/google/uuid"
)

type PlaceProduct struct {
	ID         string
	FromDate   cbigint
	Place_ID   string
	Product_ID string
}

func NewPlaceProduct() PlaceProduct {
	return PlaceProduct{
		ID:         uuid.NewString(),
		FromDate:   cbigint{Val: time.Now().Unix()},
		Place_ID:   "",
		Product_ID: "",
	}
}
