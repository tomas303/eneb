package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_placespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Place, error) {
		place := data.NewPlace()
		err := row.Scan(&place.ID, &place.Name, &place.CircuitBreakerCurrent)
		if err != nil {
			return nil, err
		}
		return &place, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`select id, name, circuitbreakercurrent
		from places
		where (name, id) < (?, ?)
		order by name desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`select id, name, circuitbreakercurrent
		from places
		where (name, id) > (?, ?)
		order by name, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/places/page/prev",
		func(c *gin.Context) {
			name := ctxQParamStr(c, "name")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*name, *id, *limit})
		})

	r.GET("/places/page/next",
		func(c *gin.Context) {
			name := ctxQParamStr(c, "name")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*name, *id, *limit})
		})

}
