package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_lastenergies(r *gin.Engine, db *sql.DB) {

	getcountparam := MakeGetQueryParAsInt("count", 10, abortWith(http.StatusBadRequest))

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
	r.GET("/lastenergies", MakeHandlerGetMany[*data.Energy]([]ParamGetterFunc{getcountparam}, cmdSelectMany))
}
