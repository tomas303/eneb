package data

import "github.com/google/uuid"

type Place struct {
	ID                    string
	Name                  string
	CircuitBreakerCurrent int
}

func NewPlace() Place {
	return Place{
		ID:                    uuid.NewString(),
		Name:                  "???",
		CircuitBreakerCurrent: 0,
	}
}
