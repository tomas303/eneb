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
			created INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tag TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS energiestags (
			enerygy_id TEXT,
			tag_id INTEGER,
			PRIMARY KEY(enerygy_id, tag_id),
			FOREIGN KEY(enerygy_id) REFERENCES energies(id),
			FOREIGN KEY(tag_id) REFERENCES tags(id)
			);`,
	}

	for _, stmt := range sqlStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement %q: %v", stmt, err)
		}
	}

	return nil
}
