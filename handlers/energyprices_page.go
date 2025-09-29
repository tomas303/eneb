package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energypricespaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.EnergyPrice, error) {
		en := data.NewEnergyPrice()
		err := row.Scan(&en.ID, &en.Kind, &en.FromDate.Val, &en.Price_ID, &en.Place_ID)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`select id, kind, fromdate, price_id, place_id
		from energyprices
		where (fromdate, id) < (?, ?)
		order by fromdate desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`select id, kind, fromdate, price_id, place_id
		from energyprices
		where (fromdate, id) > (?, ?)
		order by fromdate, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/energyprices/page/prev",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*fromdate, *id, *limit})
		})

	r.GET("/energyprices/page/next",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*fromdate, *id, *limit})
		})

}
