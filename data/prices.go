package data

import (
	"time"

	"github.com/google/uuid"
)

type Price struct {
	ID         string
	Value      int
	FromDate   cdate
	Product_ID string
	PriceType  int
	EnergyKind int
}

func NewPrice() Price {
	return Price{
		ID:         uuid.NewString(),
		Value:      0,
		FromDate:   cdate{time.Now()},
		Product_ID: "???",
		PriceType:  0,
		EnergyKind: 0,
	}
}
