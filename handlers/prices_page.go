package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_pricespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Price, error) {
		price := data.NewPrice()
		err := row.Scan(&price.ID, &price.Value, &price.FromDate, &price.Provider_ID, &price.PriceType)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.Price](db,
		`SELECT id, value, fromdate, provider_id, pricetype
		FROM prices 
		WHERE (fromdate, id) < (?, ?)
		ORDER BY fromdate DESC, id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany[*data.Price](cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.Price](db,
		`SELECT id, value, fromdate, provider_id, pricetype
		FROM prices 
		WHERE (fromdate, id) > (?, ?)
		ORDER BY fromdate, id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany[*data.Price](cmdSelectAfter)

	r.GET("/prices/page/prev",
		func(c *gin.Context) {
			fromdate := ctxQParamStr(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*fromdate, *id, *limit})
		})

	r.GET("/prices/page/next",
		func(c *gin.Context) {
			fromdate := ctxQParamStr(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*fromdate, *id, *limit})
		})

}
