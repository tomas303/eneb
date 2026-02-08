package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_productpricespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.ProductPrice, error) {
		price := data.NewProductPrice()
		err := row.Scan(&price.ID, &price.Product_ID, &price.FromDate.Val, &price.Value)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, product_id, fromdate, value
		FROM productprices
		WHERE (fromdate, id) < (?, ?)
		ORDER BY fromdate DESC, id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, product_id, fromdate, value
		FROM productprices
		WHERE (fromdate, id) > (?, ?)
		ORDER BY fromdate, id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/productprices/page/prev",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*fromdate, *id, *limit})
		})

	r.GET("/productprices/page/next",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*fromdate, *id, *limit})
		})

}
