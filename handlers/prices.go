package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_prices(r *gin.Engine, db *sql.DB) {
	r.GET("/prices", func(c *gin.Context) {
		rows, err := db.Query("SELECT value, fromdate, provider_id, energykind FROM prices")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var prices []data.Price
		for rows.Next() {
			var price data.Price
			if err := rows.Scan(&price.Value, &price.FromDate, &price.Provider_ID, &price.EnergyKind); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			prices = append(prices, price)
		}

		c.JSON(http.StatusOK, prices)
	})

	r.POST("/prices", func(c *gin.Context) {
		var price data.Price
		if err := c.ShouldBindJSON(&price); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT OR REPLACE INTO prices (value, fromdate, provider_id, energykind) VALUES (?, ?, ?, ?)", price.Value, price.FromDate, price.Provider_ID, price.EnergyKind)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, price)
	})
}
