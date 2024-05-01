package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"database/sql"
	"log"

	"eneb/config"
	"eneb/data"

	_ "github.com/glebarez/go-sqlite"
)

func main() {

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
	r.Use(gin.Recovery())
	r.Use(errHandler)
	r.Use(corsMiddleware())
	r.GET("/energies",
		func(c *gin.Context) {
			getEnergies(c, db)
		})
	r.GET("/energies/:id",
		func(c *gin.Context) {
			getEnergy(c, db)
		})
	r.POST("/energies",
		func(c *gin.Context) {
			postEnergy(c, db)
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

func getEnergies(c *gin.Context, db *sql.DB) {
	rows, err := data.LoadEnergies(db)
	if err != nil {
		c.Set("error", err)
		return
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
