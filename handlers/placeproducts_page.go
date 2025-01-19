package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_placeproductspaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.PlaceProduct, error) {
		pp := data.NewPlaceProduct()
		err := row.Scan(&pp.ID, &pp.FromDate, &pp.Place_ID, &pp.Product_ID)
		if err != nil {
			return nil, err
		}
		return &pp, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany[*data.PlaceProduct](db,
		`SELECT id, fromdate, place_id, product_id
		FROM placeproducts 
		WHERE (fromdate, id) < (?, ?)
		ORDER BY fromdate DESC, id DESC LIMIT ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany[*data.PlaceProduct](cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany[*data.PlaceProduct](db,
		`SELECT id, fromdate, place_id, product_id
		FROM placeproducts 
		WHERE (fromdate, id) > (?, ?)
		ORDER BY fromdate, id LIMIT ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany[*data.PlaceProduct](cmdSelectAfter)

	r.GET("/placeproducts/page/prev",
		func(c *gin.Context) {
			fromdate := ctxQParamStr(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*fromdate, *id, *limit})
		})

	r.GET("/placeproducts/page/next",
		func(c *gin.Context) {
			fromdate := ctxQParamStr(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*fromdate, *id, *limit})
		})

}
