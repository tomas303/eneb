package handlers

import (
	"github.com/gin-gonic/gin"
)

func Reg_root(r *gin.Engine) {
	r.GET("/", rootHandler)
}

func rootHandler(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Welcome to enef",
	})

}
