package data

import (
	"database/sql"
	"fmt"
	"log"
)

func OpenDB(dbpath string) *sql.DB {
	db, err := sql.Open("sqlite", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	err = initDB(db)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	return db
}

func initDB(db *sql.DB) error {
	sqlStatements := []string{
		`CREATE TABLE IF NOT EXISTS energies (
			id TEXT,
			kind INTEGER,
			amount INTEGER,
			info TEXT,
			created INTEGER,
			PRIMARY KEY (id)
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ENERGIES_I_PG ON ENERGIES (created, id);`,
	}

	for _, stmt := range sqlStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement %q: %v", stmt, err)
		}
	}

	return nil
}
