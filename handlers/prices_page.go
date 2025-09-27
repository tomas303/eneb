package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_pricespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Price, error) {
		price := data.NewPrice()
		err := row.Scan(&price.ID, &price.Value, &price.EnergyKind, &price.PriceType, &price.Provider_ID, &price.Name)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.Price](db,
		`SELECT id, value, energykind, pricetype, provider_id, name
		FROM prices 
		WHERE id < ?
		ORDER BY id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany[*data.Price](cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.Price](db,
		`SELECT id, value, energykind, pricetype, provider_id, name
		FROM prices 
		WHERE id > ?
		ORDER BY id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany[*data.Price](cmdSelectAfter)

	r.GET("/prices/page/prev",
		func(c *gin.Context) {
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*id, *limit})
		})

	r.GET("/prices/page/next",
		func(c *gin.Context) {
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*id, *limit})
		})

}
