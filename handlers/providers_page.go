package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_providerspaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Provider, error) {
		provider := data.NewProvider()
		err := row.Scan(&provider.ID, &provider.Name)
		if err != nil {
			return nil, err
		}
		return &provider, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, name
		FROM providers 
		WHERE id < ?
		ORDER BY id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, name
		FROM providers 
		WHERE id > ?
		ORDER BY id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/providers/page/prev",
		func(c *gin.Context) {
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*id, *limit})
		})

	r.GET("/providers/page/next",
		func(c *gin.Context) {
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*id, *limit})
		})

}
