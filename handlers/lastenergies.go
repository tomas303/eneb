package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_lastenergies(r *gin.Engine, db *sql.DB) {
	cmdSelectMany, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		order by created desc limit ?`,
		true,
		func(row data.RowScanner) (*data.Energy, error) {
			en := data.NewEnergy()
			err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created.Val)
			if err != nil {
				return nil, err
			}
			return &en, nil
		})
	if err != nil {
		panic(err)
	}
	handler := MakeHandlerGetMany[*data.Energy](cmdSelectMany)
	r.GET("/lastenergies",
		func(c *gin.Context) {
			count := ctxQParamInt(c, "count")
			handler(c, []any{count})
		})
}
