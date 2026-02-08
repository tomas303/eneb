package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_placeproducts(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.PlaceProduct, error) {
		placeProduct := data.NewPlaceProduct()
		err := row.Scan(
			&placeProduct.ID,
			&placeProduct.FromDate.Val,
			&placeProduct.Place_ID,
			&placeProduct.Product_ID,
		)
		if err != nil {
			return nil, err
		}
		return &placeProduct, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, fromdate, place_id, product_id
		FROM placeproducts
		ORDER BY fromdate DESC, id`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/placeproducts",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"INSERT OR REPLACE INTO placeproducts(id, fromdate, place_id, product_id) VALUES(?,?,?,?)",
		func(placeProduct *data.PlaceProduct) []any {
			return []any{placeProduct.ID, placeProduct.FromDate.Val, placeProduct.Place_ID, placeProduct.Product_ID}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/placeproducts",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
