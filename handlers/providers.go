package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_providers(r *gin.Engine, db *sql.DB) {
	r.GET("/providers", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name FROM providers")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var providers []data.Provider
		for rows.Next() {
			var provider data.Provider
			if err := rows.Scan(&provider.ID, &provider.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			providers = append(providers, provider)
		}

		c.JSON(http.StatusOK, providers)
	})

	r.POST("/providers", func(c *gin.Context) {
		var provider data.Provider
		if err := c.ShouldBindJSON(&provider); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT OR REPLACE INTO providers (id, name) VALUES (?, ?)", provider.ID, provider.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, provider)
	})
}
