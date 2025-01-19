package data

import (
	"time"

	"github.com/google/uuid"
)

type PlaceProduct struct {
	ID         string
	FromDate   cdate
	Place_ID   string
	Product_ID string
}

func NewPlaceProduct() PlaceProduct {
	return PlaceProduct{
		ID:         uuid.NewString(),
		FromDate:   cdate{time.Now()},
		Place_ID:   "???",
		Product_ID: "???",
	}
}
