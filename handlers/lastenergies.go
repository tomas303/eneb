package handlers

import (
	"database/sql"
	"eneb/data"
	"eneb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_lastenergies(r *gin.Engine, db *sql.DB) {
	fcountparam := utils.QParamAsInt("count", 10, abortWith(http.StatusBadRequest))
	r.GET("/lastenergies",
		func(c *gin.Context) {
			count := fcountparam(c)
			getLastEnergies(c, db, count)
		})
}

func getLastEnergies(c *gin.Context, db *sql.DB, count int) {
	var rows *[]data.Energy
	var err error
	rows, err = data.LoadLastEnergies(db, count)
	if err != nil {
		c.AbortWithError(400, err)
	}
	c.IndentedJSON(http.StatusOK, rows)
}
