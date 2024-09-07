package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Reg_common(r *gin.Engine) {
	r.Use(gin.Recovery())
	r.Use(errHandler)
	r.Use(corsMiddleware())
}

func errHandler(c *gin.Context) {
	defer func() {
		if err, exists := c.Get("error"); exists {
			if err, ok := err.(error); ok {
				errResponse := errorResponse{
					Error: err.Error(),
				}
				c.IndentedJSON(http.StatusBadRequest, errResponse)
			} else {
				errResponse := errorResponse{
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