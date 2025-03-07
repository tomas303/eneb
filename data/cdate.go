package data

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type cdate struct {
	time.Time
}

func (cd cdate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, cd.Format("20060102"))), nil
}

func (cd *cdate) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(`"20060102"`, string(data))
	if err != nil {
		return err
	}
	cd.Time = parsedTime
	return nil
}

func (cd cdate) Value() (driver.Value, error) {
	return cd.Format("20060102"), nil
}

func (cd *cdate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		cd.Time = v
	case string:
		parsedTime, err := time.Parse("20060102", v)
		if err != nil {
			return err
		}
		cd.Time = parsedTime
	default:
		return fmt.Errorf("cannot scan type %T into cdate", value)
	}
	return nil
}
