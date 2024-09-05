package main

import (
	"fmt"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"database/sql"
	"log"

	"eneb/config"
	"eneb/data"
	"eneb/utils"

	_ "github.com/glebarez/go-sqlite"
)

func main() {

	gin.EnableJsonDecoderUseNumber()
	// configuration
	cfgFile := "config.toml"
	cfg, err := config.New(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Storepath:\t", cfg.Data.Storepath)
	log.Println("Port:\t", cfg.Server.Port)
	log.Println("Release:\t", cfg.Server.Release)

	if cfg.Server.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	// database
	db := data.Open(filepath.Join(cfg.Data.Storepath, "energies.db"))
	defer db.Close()

	// routes
	r := gin.New()

	fabort := func(c *gin.Context, err error) { c.AbortWithError(400, err) }
	fprevparam := utils.QParamAsInt("prev", 0, fabort)
	fnextparam := utils.QParamAsInt("next", 10, fabort)
	fpinparam := utils.QParamAsInt64("pin", math.MaxInt64, fabort)
	fcountparam := utils.QParamAsInt("count", 10, fabort)

	r.Use(gin.Recovery())
	r.Use(errHandler)
	r.Use(corsMiddleware())
	r.GET("/energies",
		func(c *gin.Context) {
			prev := fprevparam(c)
			next := fnextparam(c)
			pin := fpinparam(c)
			if prev != 0 && next != 0 {
				c.AbortWithError(400, paramErr{message: "cannot specify both prev and next parameter"})
			}
			if (prev != 0 || next != 0) && pin == 0 {
				c.AbortWithError(400, paramErr{message: "for prev or next parameter the pin parameter is mandatory"})
			}
			getEnergies(c, db, prev, next, pin)
		})
	r.GET("/energies/:id",
		func(c *gin.Context) {
			getEnergy(c, db)
		})
	r.POST("/energies",
		func(c *gin.Context) {
			postEnergy(c, db)
		})
	r.GET("/lastenergies",
		func(c *gin.Context) {
			count := fcountparam(c)
			getLastEnergies(c, db, count)
		})

	// start
	r.Run(fmt.Sprintf("localhost:%d", cfg.Server.Port))
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errHandler(c *gin.Context) {
	defer func() {
		if err, exists := c.Get("error"); exists {
			if err, ok := err.(error); ok {
				errResponse := ErrorResponse{
					Error: err.Error(),
				}
				c.IndentedJSON(http.StatusBadRequest, errResponse)
			} else {
				errResponse := ErrorResponse{
					Error: "doesn't support error interface",
				}
				c.IndentedJSON(http.StatusInternalServerError, errResponse)
			}
		}
	}()
	c.Next()
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin with the specified methods
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
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

func getLastEnergies(c *gin.Context, db *sql.DB, count int) {
	var rows *[]data.Energy
	var err error
	rows, err = data.LoadLastEnergies(db, count)
	if err != nil {
		c.AbortWithError(400, err)
	}
	c.IndentedJSON(http.StatusOK, rows)
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

type paramErr struct {
	message string
}

func (e paramErr) Error() string {
	return e.message
}
