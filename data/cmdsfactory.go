package data

import (
	"database/sql"
)

type RowScanner interface {
	Scan(row ...any) error
}

type DataScannerFunc[T any] func(row RowScanner) (T, error)
type DataSlicerFunc[T any] func(T) []any

type DataCmdSelectOneFunc[T any] func(any) (T, error)
type DataCmdSelectManyFunc[T any] func([]any) ([]T, error)
type DataCmdSaveOneFunc[T any] func(T) (T, error)

func MakeDataCmdSelectOne[T any](db *sql.DB, sqlText string, scanner DataScannerFunc[T]) (DataCmdSelectOneFunc[T], error) {
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return nil, err
	}
	return func(id any) (T, error) {
		row := stmt.QueryRow(id)
		x, err := scanner(row)
		return x, err
	}, nil
}

func MakeDataCmdSelectMany[T any](db *sql.DB, sqlText string, scanner DataScannerFunc[T]) (DataCmdSelectManyFunc[T], error) {
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return nil, err
	}
	return func(params []any) ([]T, error) {
		rows, err := stmt.Query(params...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		data := make([]T, 0, 100)
		for rows.Next() {
			x, err := scanner(rows)
			if err != nil {
				return nil, err
			}
			data = append(data, x)
		}
		return data, nil
	}, nil
}

func MakeDataCmdSaveOne[T any](db *sql.DB, sqlText string, slicer DataSlicerFunc[T]) (DataCmdSaveOneFunc[T], error) {
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		return nil, err
	}
	return func(en T) (T, error) {
		_, err = stmt.Exec(slicer(en)...)
		if err != nil {
			return en, err
		}
		return en, nil
	}, nil
}
