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

type SQLCommand struct {
	Statement string
	ShouldRun func(db *sql.DB) bool
}

func columnExists(db *sql.DB, tableName, columnName string) bool {
	var exists int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM pragma_table_info('%s') WHERE name='%s';`, tableName, columnName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("error checking for column '%s' in table '%s': %v", columnName, tableName, err)
		return false
	}
	return exists > 0
}

func tableExists(db *sql.DB, tableName string) bool {
	var exists int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s';`, tableName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("error checking for table '%s': %v", tableName, err)
		return false
	}
	return exists > 0
}

func initDB(db *sql.DB) error {
	sqlCommands := []SQLCommand{
		{
			Statement: `CREATE TABLE places (
				id TEXT,
				name TEXT,
				circuitbreakercurrent INTEGER,
				PRIMARY KEY (id)
			);`,
			ShouldRun: func(db *sql.DB) bool {
				return !tableExists(db, "places")
			},
		},
		{
			Statement: `CREATE UNIQUE INDEX IF NOT EXISTS PLACES_I_PG ON PLACES (name, id);`,
			ShouldRun: nil,
		},
		{
			Statement: `CREATE TABLE energies (
				id TEXT,
				kind INTEGER,
				amount INTEGER,
				info TEXT,
				created INTEGER,
				PRIMARY KEY (id)
			);`,
			ShouldRun: func(db *sql.DB) bool {
				return !tableExists(db, "energies")
			},
		},
		{
			Statement: `CREATE UNIQUE INDEX IF NOT EXISTS ENERGIES_I_PG ON ENERGIES (created, id);`,
			ShouldRun: nil,
		},
		{
			Statement: `ALTER TABLE energies ADD COLUMN place_id TEXT REFERENCES places(id);`,
			ShouldRun: func(db *sql.DB) bool {
				return !columnExists(db, "energies", "place_id")
			},
		},
		{
			Statement: `CREATE TABLE providers (
				id TEXT,
				name TEXT,
				PRIMARY KEY (id)
			);`,
			ShouldRun: func(db *sql.DB) bool {
				return !tableExists(db, "providers")
			},
		},
		{
			Statement: `CREATE TABLE prices (
				id TEXT,
				value INTEGER,
				energykind INTEGER,
				pricetype INTEGER,
				provider_id TEXT,
				name TEXT,
				PRIMARY KEY (id),
				FOREIGN KEY (provider_id) REFERENCES providers(id)
			);`,
			ShouldRun: func(db *sql.DB) bool {
				return !tableExists(db, "prices")
			},
		},
	}

	for _, cmd := range sqlCommands {
		if cmd.ShouldRun == nil || cmd.ShouldRun(db) {
			_, err := db.Exec(cmd.Statement)
			if err != nil {
				return fmt.Errorf("error executing statement %q: %v", cmd.Statement, err)
			}
		}
	}

	return nil
}
