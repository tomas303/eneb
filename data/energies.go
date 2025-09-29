package data

import (
	"time"

	"github.com/google/uuid"
)

type Energy struct {
	ID       string
	Kind     int
	Amount   int
	Info     string
	Created  cbigint
	Place_ID string
}

func NewEnergy() Energy {
	return Energy{
		ID:       uuid.NewString(),
		Kind:     0,
		Amount:   0,
		Info:     "",
		Created:  cbigint{Val: time.Now().Unix()},
		Place_ID: "",
	}
}
