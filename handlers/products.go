package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_products(r *gin.Engine, db *sql.DB) {

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

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, energykind, pricetype, provider_id, name
		FROM products 
		ORDER BY name, id`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/products",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"INSERT OR REPLACE INTO products(id, energykind, pricetype, provider_id, name) VALUES(?,?,?,?,?)",
		func(product *data.Product) []any {
			return []any{product.ID, product.EnergyKind, product.PriceType, product.Provider_ID, product.Name}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/products",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
