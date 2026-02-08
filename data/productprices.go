package data

import (
	"time"

	"github.com/google/uuid"
)

type ProductPrice struct {
	ID         string
	Product_ID string
	FromDate   cbigint
	Value      int
}

func NewProductPrice() ProductPrice {
	return ProductPrice{
		ID:         uuid.NewString(),
		Product_ID: "",
		FromDate:   cbigint{Val: time.Now().Unix()},
		Value:      0,
	}
}
