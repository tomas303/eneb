package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_productspaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Product, error) {
		product := data.NewProduct()
		err := row.Scan(
			&product.ID,
			&product.EnergyKind,
			&product.PriceType,
			&product.Provider_ID,
			&product.Name,
		)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, energykind, pricetype, provider_id, name
		FROM products
		WHERE (name, id) < (?, ?)
		ORDER BY name DESC, id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, energykind, pricetype, provider_id, name
		FROM products
		WHERE (name, id) > (?, ?)
		ORDER BY name, id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/products/page/prev",
		func(c *gin.Context) {
			name := ctxQParamStr(c, "name")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*name, *id, *limit})
		})

	r.GET("/products/page/next",
		func(c *gin.Context) {
			name := ctxQParamStr(c, "name")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*name, *id, *limit})
		})

}
