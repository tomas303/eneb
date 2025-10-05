package data

type GasPriceSerie struct {
	Place_ID         string
	SliceStart       cbigint
	SliceEnd         cbigint
	AmountMwh        float64
	Months           float64
	UnregulatedPrice float64
	RegulatedPrice   float64
	TotalPrice       float64
}

func NewGasPriceSerie() GasPriceSerie {
	return GasPriceSerie{
		Place_ID:         "",
		SliceStart:       cbigint{},
		SliceEnd:         cbigint{},
		AmountMwh:        0,
		Months:           0,
		UnregulatedPrice: 0,
		RegulatedPrice:   0,
		TotalPrice:       0,
	}
}
