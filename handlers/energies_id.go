package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energiesid(r *gin.Engine, db *sql.DB) {

	cmdSelectOne, err := data.MakeDataCmdSelectOne[*data.Energy](db,
		"select id, kind, amount, info, created from energies where id = ?",
		func(row data.RowScanner) (*data.Energy, error) {
			en := data.NewEnergy()
			err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created)
			if err != nil {
				return nil, err
			}
			return &en, nil
		})
	if err != nil {
		panic(err)
	}

	handler := MakeHandlerGetOne[*data.Energy](cmdSelectOne)
	r.GET("/energies/:id",
		func(c *gin.Context) {
			id := ctxPParam(c, "id")
			handler(c, []any{*id})
		})
}
