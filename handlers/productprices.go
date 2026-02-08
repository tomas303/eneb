package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_productprices(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.ProductPrice, error) {
		price := data.NewProductPrice()
		err := row.Scan(
			&price.ID,
			&price.Product_ID,
			&price.FromDate.Val,
			&price.Value,
		)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, product_id, fromdate, value
		FROM productprices 
		ORDER BY fromdate DESC, id`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/productprices",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"INSERT OR REPLACE INTO productprices(id, product_id, fromdate, value) VALUES(?,?,?,?)",
		func(price *data.ProductPrice) []any {
			return []any{price.ID, price.Product_ID, price.FromDate.Val, price.Value}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/productprices",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
