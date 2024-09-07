package main

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"log"

	"eneb/config"
	"eneb/data"
	"eneb/handlers"

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

	handlers.Reg_common(r)
	handlers.Reg_energies(r, db)
	handlers.Reg_energiesid(r, db)
	handlers.Reg_lastenergies(r, db)

	// start
	r.Run(fmt.Sprintf("localhost:%d", cfg.Server.Port))
}
