package data

import (
	"database/sql"
	"gobackend/utils"
	"log"
)

type Data struct {
	db *sql.DB
}

func New(dbpath string) *Data {
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	err = prepare(db)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	return &Data{db: db}
}

func (data *Data) Close() {
	data.db.Close()
}

func (data *Data) LoadEnergy(id int32) (*Energy, error) {
	rows, err := data.db.Query("select amount, info from energies where id = :id", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	energy := NewEnergy()
	if rows.Next() {
		err := rows.Scan(&energy.Amount, &energy.Info)
		if err != nil {
			return nil, err
		}
		energy.ID = id
	}
	return &energy, nil
}

func (data *Data) LoadEnergies() (*[]Energy, error) {
	rows, err := data.db.Query("select amount, info from energies order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var energies []Energy
	for rows.Next() {
		energy := NewEnergy()
		err := rows.Scan(&energy.ID, &energy.Amount, &energy.Info)
		if err != nil {
			return nil, err
		}
		energies = append(energies, energy)
	}
	return &energies, nil
}

func (data *Data) LoadEnergies2() *utils.Iterator[Energy] {
	rows, err := data.db.Query("select amount, info from energies order by id")
	if err != nil {
		return nil
	}
	return utils.NewIterator[Energy](
		func(channel utils.IteratorChannel[Energy]) {
			defer rows.Close()
			for {
				if rows.Next() {
					energy := NewEnergy()
					err := rows.Scan(&energy.ID, &energy.Amount, &energy.Info)
					var x utils.Result[Energy, error]
					if err != nil {
						x.Err = err
					} else {
						x.Value = energy
					}
					channel <- x
				} else {
					break
				}
			}
		})
}

func prepare(db *sql.DB) error {
	// Create the table if it doesn't exist.
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS energies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount INTEGER,
			info TEXT
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tag TEXT
		)
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS energiestags (
			enerygy_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY(enerygy_id, tag_id),
			FOREIGN KEY(enerygy_id) REFERENCES energies(id),
			FOREIGN KEY(tag_id) REFERENCES tags(id)
		)
	`)
	if err != nil {
		return err
	}
	return nil
}
