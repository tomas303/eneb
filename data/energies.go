package data

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Energy struct {
	ID      string
	Kind    int
	Amount  int
	Info    string
	Created cbigint
}

func NewEnergy() Energy {
	return Energy{
		ID:      uuid.NewString(),
		Kind:    0,
		Amount:  0,
		Info:    "",
		Created: cbigint{Val: time.Now().Unix()},
	}
}

type cbigint struct {
	Val int64
}

func (cbi *cbigint) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	x, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	cbi.Val = x
	return nil
}

func (cbi cbigint) MarshalJSON() ([]byte, error) {
	value := strconv.FormatInt(cbi.Val, 10)
	return json.Marshal(value)
}
