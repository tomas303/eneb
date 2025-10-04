package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energiespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Energy, error) {
		en := data.NewEnergy()
		err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created.Val, &en.Place_ID)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`select id, kind, amount, info, created, place_id
		from energies 
		where (created, id) < (?, ?)
		order by created desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`select id, kind, amount, info, created, place_id
		from energies 
		where (created, id) > (?, ?)
		order by created, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/energies/page/prev",
		func(c *gin.Context) {
			created := ctxQParamInt(c, "created")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*created, *id, *limit})
		})

	r.GET("/energies/page/next",
		func(c *gin.Context) {
			created := ctxQParamInt(c, "created")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*created, *id, *limit})
		})

}
