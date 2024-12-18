package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energies(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Energy, error) {
		en := data.NewEnergy()
		err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created.Val)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		order by created, id`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany[*data.Energy](cmdSelect)

	r.GET("/energies",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne[*data.Energy](db,
		"insert or replace into energies(id, kind, amount, info, created) VALUES(?,?,?,?,?)",
		func(en *data.Energy) []any {
			return []any{en.ID, en.Kind, en.Amount, en.Info, en.Created.Val}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne[*data.Energy](cmdSave)

	r.POST("/energies",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
