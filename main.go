package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"embed"
	"log"

	"eneb/config"
	"eneb/data"
	"eneb/handlers"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed dist/*
var distFS embed.FS

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
	db := data.OpenDB(filepath.Join(cfg.Data.Storepath, "energies.db"))
	defer db.Close()

	// routes
	r := gin.New()

	// embedded static files
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(distFS, "dist/*.html")))
	assetsFS, _ := fs.Sub(distFS, "dist/assets")
	r.StaticFS("/assets", http.FS(assetsFS))

	handlers.Reg_common(r)
	handlers.Reg_root(r)
	handlers.Reg_energies(r, db)
	handlers.Reg_energiesid(r, db)
	handlers.Reg_lastenergies(r, db)

	// start
	r.Run(fmt.Sprintf("localhost:%d", cfg.Server.Port))
}
