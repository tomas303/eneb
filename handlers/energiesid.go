package handlers

import (
	"database/sql"
	"eneb/data"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Reg_energiesid(r *gin.Engine, db *sql.DB) {
	r.GET("/energies/:id",
		func(c *gin.Context) {
			getEnergy(c, db)
		})
}

func getEnergy(c *gin.Context, db *sql.DB) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Set("error", err)
		return
	}
	energy, err := data.LoadEnergy(db, id)
	if err != nil {
		c.Set("error", err)
		return
	}
	c.IndentedJSON(http.StatusOK, energy)
}
