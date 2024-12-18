package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energiespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Energy, error) {
		en := data.NewEnergy()
		err := row.Scan(&en.ID, &en.Kind, &en.Amount, &en.Info, &en.Created.Val)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where (created, id) < (?, ?)
		order by created desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany[*data.Energy](cmdSelectBefore)

	cmdSelectBeforeIncluded, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where (created, id) <= (?, ?)
		order by created desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeBoundaryHandler := MakeHandlerGetMany[*data.Energy](cmdSelectBeforeIncluded)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where (created, id) > (?, ?)
		order by created, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany[*data.Energy](cmdSelectAfter)

	cmdSelectAfterIncluded, err := data.MakeDataCmdSelectMany[*data.Energy](db,
		`select id, kind, amount, info, created 
		from energies 
		where (created, id) >= (?, ?)
		order by created, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterBoundaryHandler := MakeHandlerGetMany[*data.Energy](cmdSelectAfterIncluded)

	r.GET("/energies/page/prev",
		func(c *gin.Context) {
			created := ctxQParamInt(c, "created")
			id := ctxQParamStr(c, "id")
			included := ctxQParamBool(c, "included")
			limit := ctxQParamInt(c, "limit")
			if included != nil && *included {
				beforeBoundaryHandler(c, []any{*created, *id, *limit})
			} else {
				beforeHandler(c, []any{*created, *id, *limit})
			}
		})

	r.GET("/energies/page/next",
		func(c *gin.Context) {
			created := ctxQParamInt(c, "created")
			id := ctxQParamStr(c, "id")
			included := ctxQParamBool(c, "included")
			limit := ctxQParamInt(c, "limit")
			if included != nil && *included {
				afterBoundaryHandler(c, []any{*created, *id, *limit})
			} else {
				afterHandler(c, []any{*created, *id, *limit})
			}
		})

}
