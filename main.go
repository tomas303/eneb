package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"database/sql"
	"log"

	"gobackend/config"
	"gobackend/data"

	_ "github.com/glebarez/go-sqlite"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {

	// configuration
	cfgFile := "config1.toml"
	// var cfg Config
	// err := cfg.ReadFromFile(configFile)
	cfg, err := config.New(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Storepath:\t", cfg.Data.Storepath)
	log.Println("Port:\t", cfg.Server.Port)

	// database
	data := data.New(filepath.Join(cfg.Data.Storepath, "energies.db"))
	defer data.Close()

	// routes
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("data", data)
		c.Next()
	})
	// /energies   /tags         ....  energy will append tag list, engine will transfer it to ids
	// importang get paremters .... list of tags to have .. either all or any
	//router.POST("/db/", postAlbums)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.GET("/test", getTest)

	router.Run(fmt.Sprintf("localhost:%d", cfg.Server.Port))
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

func getTest(c *gin.Context) {

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}
	// Type assert db back to *sql.DB.
	dbConn, ok := db.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to type assert database connection"})
		return
	}

	var todos []Todo
	rows, err := dbConn.Query("SELECT id, task FROM todos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, todos)
}
