package handlers

import (
	"database/sql"
	"eneb/data"

	"github.com/gin-gonic/gin"
)

func Reg_energyprices(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.EnergyPrice, error) {
		price := data.NewEnergyPrice()
		err := row.Scan(
			&price.ID,
			&price.FromDate.Val,
			&price.Price_ID,
			&price.Place_ID,
		)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany(db,
		`SELECT id, fromdate, price_id, place_id
		FROM energyprices
		ORDER BY fromdate DESC, id`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	getHandler := MakeHandlerGetMany(cmdSelect)

	r.GET("/energyprices",
		func(c *gin.Context) {
			getHandler(c, []any{})
		})

	cmdSave, err := data.MakeDataCmdSaveOne(db,
		"insert or replace into energyprices(id, fromdate, price_id, place_id) VALUES(?,?,?,?)",
		func(price *data.EnergyPrice) []any {
			return []any{price.ID, price.FromDate.Val, price.Price_ID, price.Place_ID}
		})
	if err != nil {
		panic(err)
	}
	postHandler := MakeHandlerPostOne(cmdSave)

	r.POST("/energyprices",
		func(c *gin.Context) {
			postHandler(c, []any{})
		})
}
