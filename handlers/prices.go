package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_prices(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Price, error) {
		price := data.NewPrice()
		err := row.Scan(
			&price.ID,
			&price.Value,
			&price.EnergyKind,
			&price.PriceType,
			&price.Provider_ID,
			&price.Name,
		)
		if err != nil {
			return nil, err
		}
		return &price, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany[*data.Price](db,
		`SELECT id, value, energykind, pricetype, provider_id, name
		FROM prices 
		ORDER BY id DESC`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	handler := MakeHandlerGetMany[*data.Price](cmdSelect)

	r.GET("/prices",
		func(c *gin.Context) {
			handler(c, nil)
		})

	r.POST("/prices", func(c *gin.Context) {
		var price data.Price
		if err := c.ShouldBindJSON(&price); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(
			`INSERT OR REPLACE INTO prices 
			(id, value, energykind, pricetype, provider_id, name) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			price.ID, price.Value, price.EnergyKind, price.PriceType, price.Provider_ID, price.Name,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, price)
	})
}
