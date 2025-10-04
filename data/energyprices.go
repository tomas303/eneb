package data

import (
	"time"

	"github.com/google/uuid"
)

type EnergyPrice struct {
	ID       string
	FromDate cbigint
	Price_ID string
	Place_ID string
}

func NewEnergyPrice() EnergyPrice {
	return EnergyPrice{
		ID:       uuid.NewString(),
		FromDate: cbigint{Val: time.Now().Unix()},
		Price_ID: "",
		Place_ID: "",
	}
}
