package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_placeproducts(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.PlaceProduct, error) {
		pp := data.NewPlaceProduct()
		err := row.Scan(&pp.ID, &pp.FromDate, &pp.Place_ID, &pp.Product_ID)
		if err != nil {
			return nil, err
		}
		return &pp, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany[*data.PlaceProduct](db,
		`SELECT id, fromdate, place_id, product_id
		FROM placeproducts 
		ORDER BY fromdate DESC, id DESC`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	handler := MakeHandlerGetMany[*data.PlaceProduct](cmdSelect)

	r.GET("/placeproducts",
		func(c *gin.Context) {
			handler(c, nil)
		})

	r.POST("/placeproducts", func(c *gin.Context) {
		var pp data.PlaceProduct
		if err := c.ShouldBindJSON(&pp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT OR REPLACE INTO placeproducts (id, fromdate, place_id, product_id) VALUES (?, ?, ?, ?)", pp.ID, pp.FromDate, pp.Place_ID, pp.Product_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, pp)
	})
}
