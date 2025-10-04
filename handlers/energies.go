package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energies(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Energy, error) {
		en := data.NewEnergy()
		err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created.Val, &en.Place_ID)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`select id, kind, amount, info, created, place_ 
		from energies 
		order by created, id`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/energies",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"insert or replace into energies(id, kind, amount, info, created, place_id) VALUES(?,?,?,?,?,?)",
		func(en *data.Energy) []any {
			return []any{en.ID, en.Kind, en.Amount, en.Info, en.Created.Val, en.Place_ID}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/energies",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
