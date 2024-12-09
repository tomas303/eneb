package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_energies(r *gin.Engine, db *sql.DB) {

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where created < ?
		order by created desc limit ?`,
		false,
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
	beforeHandler := MakeHandlerGetMany[*data.Energy](cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where created > ?
		order by created limit ?`,
		false,
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
	afterHandler := MakeHandlerGetMany[*data.Energy](cmdSelectAfter)
	r.GET("/energies",
		func(c *gin.Context) {
			prev := ctxQParamInt(c, "prev")
			next := ctxQParamInt(c, "next")
			pin := ctxQParamInt(c, "pin")
			if prev != nil && next != nil {
				c.AbortWithError(http.StatusBadRequest, paramErr{message: "cannot specify both prev and next parameter"})
				return
			}
			if (prev != nil || next != nil) && pin == nil {
				c.AbortWithError(http.StatusBadRequest, paramErr{message: "for prev or next parameter the pin parameter is mandatory"})
				return
			}
			if prev != nil {
				beforeHandler(c, []any{*pin, *prev})
			} else if next != nil {
				afterHandler(c, []any{*pin, *next})
			} else {
				c.AbortWithError(400, paramErr{message: "nor prev nor next parameter specified"})
			}
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
