package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_places(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Place, error) {
		place := data.NewPlace()
		err := row.Scan(&place.ID, &place.Name, &place.CircuitBreakerCurrent)
		if err != nil {
			return nil, err
		}
		return &place, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany[*data.Place](db,
		`select id, name, circuitbreakercurrent
		from places
		order by name, id`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany[*data.Place](cmdSelect)

	r.GET("/places",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne[*data.Place](db,
		"insert or replace into places(id, name, circuitbreakercurrent) VALUES(?,?,?)",
		func(place *data.Place) []any {
			return []any{place.ID, place.Name, place.CircuitBreakerCurrent}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne[*data.Place](cmdSave)

	r.POST("/places",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
