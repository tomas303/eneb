package data

import (
	"time"

	"github.com/google/uuid"
)

type Settlement struct {
	ID         string
	Date       cbigint
	EnergyKind int
	PriceType  int
	Amount     int
	Price      int
}

func NewSettlement() Settlement {
	return Settlement{
		ID:         uuid.NewString(),
		Date:       cbigint{Val: time.Now().Unix()},
		EnergyKind: 0,
		PriceType:  0,
		Amount:     0,
		Price:      0,
	}
}
