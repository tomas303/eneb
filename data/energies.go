package data

import (
	"time"

	"github.com/google/uuid"
)

type Energy struct {
	ID      string
	Amount  int64
	Info    string
	Created int64
}

func NewEnergy() Energy {
	return Energy{
		ID:      uuid.NewString(),
		Amount:  0,
		Info:    "",
		Created: time.Now().Unix(),
	}
}
