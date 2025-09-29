package data

import (
	"time"

	"github.com/google/uuid"
)

type EnergyPrice struct {
	ID       string
	Kind     int
	FromDate cbigint
	Price_ID string
	Place_ID string
}

func NewEnergyPrice() EnergyPrice {
	return EnergyPrice{
		ID:       uuid.NewString(),
		Kind:     0,
		FromDate: cbigint{Val: time.Now().Unix()},
		Price_ID: "",
		Place_ID: "",
	}
}
