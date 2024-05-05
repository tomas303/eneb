package data

import (
	"database/sql"
	"eneb/utils"
	"fmt"
	"log"
)

func Open(dbpath string) *sql.DB {
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
	return db
}

func LoadEnergy(db *sql.DB, id int64) (*Energy, error) {
	rows, err := db.Query("select id, amount, info, created from energies where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	en := NewEnergy()
	if rows.Next() {
		err := rows.Scan(&en.ID, &en.Amount, &en.Info, &en.Created)
		if err != nil {
			return nil, err
		}
	}
	return &en, nil
}

func PostEnergy(db *sql.DB, en *Energy) (*Energy, error) {
	stmt, err := db.Prepare("insert or replace into energies(id, amount, info, created) VALUES(?,?,?,?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(en.ID, en.Amount.Val, en.Info, en.Created.Val)
	if err != nil {
		return nil, err
	}
	return en, nil
}

func LoadEnergies(db *sql.DB) (*[]Energy, error) {
	rows, err := db.Query("select id, amount, info, created from energies order by created")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var energies []Energy
	for rows.Next() {
		en := NewEnergy()
		err := rows.Scan(&en.ID, &en.Amount.Val, &en.Info, &en.Created.Val)
		if err != nil {
			return nil, err
		}
		energies = append(energies, en)
	}
	return &energies, nil
}

func LoadEnergiesAfter(db *sql.DB, pin int64, take int) (*[]Energy, error) {
	rows, err := db.Query(fmt.Sprintf(
		`select id, amount, info, created 
		from energies 
		where created > %d 
		order by created limit %d`,
		pin, take))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	energies := make([]Energy, 0, take)
	for rows.Next() {
		en := NewEnergy()
		err := rows.Scan(&en.ID, &en.Amount.Val, &en.Info, &en.Created.Val)
		if err != nil {
			return nil, err
		}
		energies = append(energies, en)
	}
	return &energies, nil
}

func LoadEnergiesBefore(db *sql.DB, pin int64, take int) (*[]Energy, error) {
	rows, err := db.Query(fmt.Sprintf(
		`select id, amount, info, created 
		from energies 
		where created < %d 
		order by created desc limit %d`,
		pin, take))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	energies := make([]Energy, 0, take)
	for rows.Next() {
		en := NewEnergy()
		err := rows.Scan(&en.ID, &en.Amount.Val, &en.Info, &en.Created.Val)
		if err != nil {
			return nil, err
		}
		energies = append([]Energy{en}, energies...)
	}
	return &energies, nil
}

func LoadEnergies2(db *sql.DB) *utils.Iterator[Energy] {
	rows, err := db.Query("select amount, info, created from energies order by id")
	if err != nil {
		return nil
	}
	return utils.NewIterator[Energy](
		func(channel utils.IteratorChannel[Energy]) {
			defer rows.Close()
			for {
				if rows.Next() {
					en := NewEnergy()
					err := rows.Scan(&en.ID, &en.Amount, &en.Info, &en.Created)
					var x utils.Result[Energy, error]
					if err != nil {
						x.Err = err
					} else {
						x.Value = en
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
			id TEXT,
			amount INTEGER,
			info TEXT,
			created INTEGER
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
			enerygy_id TEXT,
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
