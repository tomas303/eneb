package data

type Energy struct {
	ID     int32
	Amount int64
	Info   string
}

func NewEnergy() Energy {
	return Energy{
		ID:     0,
		Amount: 0,
		Info:   "",
	}
}
