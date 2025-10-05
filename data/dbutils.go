package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

func viewExists(db *sql.DB, viewName string) bool {
	var exists int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM sqlite_master WHERE type='view' AND name='%s';`, viewName)
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("error checking for view '%s': %v", viewName, err)
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
				place_id TEXT,
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
		{
			Statement: `CREATE TABLE energyprices (
				id TEXT,
				fromdate INTEGER,
				price_id TEXT,
				place_id TEXT,
				PRIMARY KEY (id),
				FOREIGN KEY (price_id) REFERENCES prices(id),
				FOREIGN KEY (place_id) REFERENCES places(id)
			);`,
			ShouldRun: func(db *sql.DB) bool {
				return !tableExists(db, "energyprices")
			},
		},
		{
			Statement: `create view v_consumptionprice as
				WITH 
				-- amount deltas(ranges)
				AmountDeltas AS (
				SELECT
				id,
				created,
				kind,
				place_id,
				LAG(created) OVER (partition by place_id,kind order by created) AS prevcreated,
				amount - LAG(amount) OVER (partition by place_id,kind order by created) AS amount_delta
				FROM
				energies
				),
				-- prices based on date they started
				PriceCoeficients AS (
				select 
				ep.fromdate, ep.place_id, p.energykind, p.value, p.pricetype 
				from
				energyprices ep
				join prices p on p.id = ep.price_id  
				),
				-- coefficient change dates inside amount deltas
				Boundaries AS (
					SELECT DISTINCT
						AD.id,
						PC.fromdate,
						PC.energykind,
						PC.place_id 
					FROM AmountDeltas AD
					JOIN PriceCoeficients PC
					ON PC.energykind = AD.kind 
					and PC.place_id  = AD.place_id 
					and PC.fromdate > AD.prevcreated
					and PC.fromdate < AD.created
				),
				-- cut points = range start, end, and coef changes
				cuts AS (
					SELECT id, kind, place_id, prevcreated AS cut_date FROM AmountDeltas
					UNION
					SELECT id, kind, place_id, created AS cut_date FROM AmountDeltas
					UNION
					SELECT id, energykind, place_id, fromdate AS cut_date FROM Boundaries
				),
				-- consecutive cuts form slices
				slices AS (
					SELECT
						id,
						kind, 
						place_id,
						cut_date AS slice_start,
						LEAD(cut_date) OVER (PARTITION BY id ORDER BY cut_date) AS slice_end
					FROM cuts
				),
				valid_slices AS (
					SELECT * FROM slices WHERE slice_start IS NOT NULL and slice_end IS NOT NULL
				),
				-- attach coefficient values valid at slice_start
				slice_coefs AS (
					SELECT
						s.id,
						s.kind,
						s.place_id,
						s.slice_start,
						s.slice_end,
						PC.pricetype,
						PC.value,
						PC.fromdate 
					FROM valid_slices s
					join PriceCoeficients PC on PC.energykind = s.kind  and PC.place_id  = s.place_id 
					and PC.fromdate = (select min(fromdate) from PriceCoeficients pcx where pcx.fromdate<=s.slice_start and pcx.energykind = s.kind and pcx.place_id = s.place_id )
				)
				select
				AD.id,
				AD.place_id, 
				AD.kind, 
				AD.created,
				AD.amount_delta,
				SC.pricetype,
				SC.value,
				SC.fromdate,
				SC.slice_start,
				SC.slice_end,
				((julianday(datetime(SC.slice_end, 'unixepoch')) - julianday(datetime(SC.slice_start, 'unixepoch'))) / 30.44) AS months_diff,
				(SC.slice_end - SC.slice_start) / (AD.created - AD.prevcreated) * AD.amount_delta AS proportional_amount
				from
				AmountDeltas AD 
				join slice_coefs SC on SC.id = AD.id
				order by AD.place_id, AD.kind, AD.created 

			;`,
			ShouldRun: func(db *sql.DB) bool {
				return !viewExists(db, "v_consumptionprice")
			},
		},
		{
			Statement: `create view v_consumptionpricegas as
				with
				pivoted AS (
				select 
					place_id, kind, 
					avg(proportional_amount) * 10.55 / 1000 as amountMwh,
					avg(months_diff) as months,
					MAX(CASE WHEN pricetype = 1 THEN value END) AS ComodityPerVolume,
					MAX(CASE WHEN pricetype = 2 THEN value END) AS DistributionPerVolume,
					MAX(CASE WHEN pricetype = 3 THEN value END) AS ComodityPerMonth,
					MAX(CASE WHEN pricetype = 4 THEN value END) AS DistributionPerMonth,
					MAX(CASE WHEN pricetype = 5 THEN value END) AS OTE,
					MAX(CASE WHEN pricetype = 11 THEN value END) AS VAT
				from v_consumptionprice
				where pricetype in (1,2,3,4,5,11)
				group by place_id, kind, slice_end
				),
				priceCalc as (
				select 
					place_id, kind, amountMwh, months, VAT,
					ROUND((amountMwh * ComodityPerVolume) + (months * ComodityPerMonth), 0) as unregulatedPrice,
					ROUND((amountMwh * DistributionPerVolume) + (months * DistributionPerMonth) + (amountMwh * OTE), 0) as regulatedPrice
				from pivoted
				)
				select 
				place_id, kind, amountMwh, months, unregulatedPrice, regulatedPrice,
				ROUND((unregulatedPrice + regulatedPrice) * VAT, 0) as totalPrice  
				from priceCalc
			;`,
			ShouldRun: func(db *sql.DB) bool {
				return !viewExists(db, "v_consumptionpricegas")
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
