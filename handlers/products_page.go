package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_productspaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Product, error) {
		product := data.NewProduct()
		err := row.Scan(&product.ID, &product.Name, &product.Provider_ID)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.Product](db,
		`SELECT id, name, provider_id
		FROM products 
		WHERE id < ?
		ORDER BY id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany[*data.Product](cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.Product](db,
		`SELECT id, name, provider_id
		FROM products 
		WHERE id > ?
		ORDER BY id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany[*data.Product](cmdSelectAfter)

	r.GET("/products/page/prev",
		func(c *gin.Context) {
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*id, *limit})
		})

	r.GET("/products/page/next",
		func(c *gin.Context) {
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*id, *limit})
		})

}
