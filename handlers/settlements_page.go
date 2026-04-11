package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_settlementspaging(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Settlement, error) {
		settlement := data.NewSettlement()
		err := row.Scan(&settlement.ID, &settlement.Date.Val, &settlement.EnergyKind, &settlement.PriceType, &settlement.Amount, &settlement.Price)
		if err != nil {
			return nil, err
		}
		return &settlement, nil
	}

	// Previous page (before cursor)
	cmdSelectBefore, err := data.MakeDataCmdSelectMany(db,
		`select id, date, energykind, pricetype, amount, price
		from settlements
		where (date, id) < (?, ?)
		order by date desc, id desc limit ?`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	beforeHandler := MakeHandlerGetMany(cmdSelectBefore)

	// Next page (after cursor)
	cmdSelectAfter, err := data.MakeDataCmdSelectMany(db,
		`select id, date, energykind, pricetype, amount, price
		from settlements
		where (date, id) > (?, ?)
		order by date, id limit ?`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	afterHandler := MakeHandlerGetMany(cmdSelectAfter)

	r.GET("/settlements/page/prev",
		func(c *gin.Context) {
			date := ctxQParamInt(c, "date")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			beforeHandler(c, []any{*date, *id, *limit})
		})

	r.GET("/settlements/page/next",
		func(c *gin.Context) {
			date := ctxQParamInt(c, "date")
			id := ctxQParamStr(c, "id")
			limit := ctxQParamInt(c, "limit")
			afterHandler(c, []any{*date, *id, *limit})
		})
}
