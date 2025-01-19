package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Reg_products(r *gin.Engine, db *sql.DB) {

	getScanner := func(row data.RowScanner) (*data.Product, error) {
		product := data.NewProduct()
		err := row.Scan(&product.ID, &product.Name, &product.ProviderID)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

	cmdSelect, err := data.MakeDataCmdSelectMany[*data.Product](db,
		`SELECT id, name, provider_id
		FROM products 
		ORDER BY name DESC, id DESC`,
		true,
		getScanner)
	if err != nil {
		panic(err)
	}
	handler := MakeHandlerGetMany[*data.Product](cmdSelect)

	r.GET("/products",
		func(c *gin.Context) {
			handler(c, nil)
		})

	r.POST("/products", func(c *gin.Context) {
		var product data.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.ID = uuid.New().String()
		_, err := db.Exec("INSERT OR REPLACE INTO products (id, name, provider_id) VALUES (?, ?, ?)", product.ID, product.Name, product.ProviderID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	})
}
