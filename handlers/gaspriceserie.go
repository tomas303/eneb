package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_gaspriceserie(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.GasPriceSerie, error) {
		en := data.NewGasPriceSerie()
		err := row.Scan(&en.Place_ID, &en.SliceStart.Val, &en.SliceEnd.Val, &en.AmountMwh, &en.Months, &en.UnregulatedPrice, &en.RegulatedPrice, &en.TotalPrice)
		if err != nil {
			return nil, err
		}
		return &en, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`select Place_ID, Slice_Start, Slice_End, AmountMwh, Months, UnregulatedPrice, RegulatedPrice, TotalPrice
		from v_consumptionpricegas
		where Slice_Start >= ? and Slice_End <= ?
		order by Slice_Start, Slice_End`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/gas-prices",
		func(c *gin.Context) {
			fromdate := ctxQParamInt(c, "fromdate")
			todate := ctxQParamInt(c, "todate")
			getHandler(c, []any{*fromdate, *todate})
		})
}
