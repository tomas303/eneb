package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_placeproductspaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.PlaceProduct, error) {
		en := data.NewPlaceProduct()
		err := row.Scan(&en.ID, &en.FromDate.Val, &en.Place_ID, &en.Product_ID)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`select id, fromdate, place_id, product_id
		from placeproducts
		where (fromdate, id) < (?, ?)
		order by fromdate desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`select id, fromdate, place_id, product_id
		from placeproducts
		where (fromdate, id) > (?, ?)
		order by fromdate, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/placeproducts/page/prev",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*fromdate, *id, *limit})
		})

	r.GET("/placeproducts/page/next",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*fromdate, *id, *limit})
		})

}
