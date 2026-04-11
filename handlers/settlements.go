package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_settlements(r *gin.Engine, db *sql.DB) {

	// GET /settlements - List all
	getScanner := func(row data.RowScanner) (*data.Settlement, error) {
		settlement := data.NewSettlement()
		err := row.Scan(&settlement.ID, &settlement.Date.Val, &settlement.EnergyKind, &settlement.PriceType, &settlement.Amount, &settlement.Price)
		if err != nil {
			return nil, err
		}
		return &settlement, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`select id, date, energykind, pricetype, amount, price
		from settlements
		order by date desc, id`,
		false,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/settlements",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	// POST /settlements - Create/Update
	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"insert or replace into settlements(id, date, energykind, pricetype, amount, price) VALUES(?,?,?,?,?,?)",
		func(settlement *data.Settlement) []any {
			return []any{settlement.ID, settlement.Date.Val, settlement.EnergyKind, settlement.PriceType, settlement.Amount, settlement.Price}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/settlements",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
