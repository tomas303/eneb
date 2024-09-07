package handlers

import (
	"database/sql"
	"eneb/data"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_energies(r *gin.Engine, db *sql.DB) {
	fprevparam := qParamAsInt("prev", 0, abortWith(http.StatusBadRequest))
	fnextparam := qParamAsInt("next", 10, abortWith(http.StatusBadRequest))
	fpinparam := qParamAsInt64("pin", math.MaxInt64, abortWith(http.StatusBadRequest))

	r.GET("/energies",
		func(c *gin.Context) {
			prev := fprevparam(c)
			next := fnextparam(c)
			pin := fpinparam(c)
			if prev != 0 && next != 0 {
				c.AbortWithError(http.StatusBadRequest, paramErr{message: "cannot specify both prev and next parameter"})
			}
			if (prev != 0 || next != 0) && pin == 0 {
				c.AbortWithError(http.StatusBadRequest, paramErr{message: "for prev or next parameter the pin parameter is mandatory"})
			}
			getEnergies(c, db, prev, next, pin)
		})

	r.POST("/energies",
		func(c *gin.Context) {
			postEnergy(c, db)
		})

}

func getEnergies(c *gin.Context, db *sql.DB, prev int, next int, pin int64) {
	var rows *[]data.Energy
	var err error
	if prev > 0 {
		rows, err = data.LoadEnergiesBefore(db, pin, prev)
	} else if next > 0 {
		rows, err = data.LoadEnergiesAfter(db, pin, next)
	} else {
		c.AbortWithError(400, paramErr{message: "nor prev nor next parameter specified"})
	}
	if err != nil {
		c.AbortWithError(400, err)
	}
	c.IndentedJSON(http.StatusOK, rows)
}

func postEnergy(c *gin.Context, db *sql.DB) {
	var energy data.Energy
	if err := c.BindJSON(&energy); err != nil {
		c.Set("error", err)
		return
	}
	_, err := data.PostEnergy(db, &energy)
	if err != nil {
		c.Set("error", err)
		return
	}
	c.Status(http.StatusCreated)
}
