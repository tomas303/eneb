package data

import "github.com/google/uuid"

type Provider struct {
	ID   string
	Name string
}

func NewProvider() Provider {
	return Provider{
		ID:   uuid.NewString(),
		Name: "???",
	}
}
