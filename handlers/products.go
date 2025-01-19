package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_products(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Product, error) {
		product := data.NewProduct()
		err := row.Scan(&product.ID, &product.Name, &product.Provider_ID)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany[*data.Product](db,
		`SELECT id, name, provider_id
		FROM products 
		ORDER BY name DESC, id DESC`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany[*data.Product](cmdSelect)

	r.GET("/products",
		func(c *gin.Context) {
			getHandler(c, nil)
		})

	cmdSave, err := data.MakeDataCmdSaveOne[*data.Product](db,
		"insert or replace into products(id, name, provider_id) VALUES(?,?,?)",
		func(product *data.Product) []any {
			return []any{product.ID, product.Name, product.Provider_ID}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne[*data.Product](cmdSave)

	r.POST("/products", func(c *gin.Context) {
		postHandler(c, []any{})
	})
}
