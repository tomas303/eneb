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
	handlers.Reg_energiespaging(r, db)
	handlers.Reg_places(r, db)
	handlers.Reg_placespaging(r, db)
	handlers.Reg_providers(r, db)
	handlers.Reg_providerspaging(r, db)
	handlers.Reg_prices(r, db)
	handlers.Reg_pricespaging(r, db)
	handlers.Reg_energypricespaging(r, db)
	handlers.Reg_energyprices(r, db)
	handlers.Reg_gaspriceserie(r, db)
	// start
	r.Run(fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port))
}
