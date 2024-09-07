package handlers

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

func abortWith(statuscode int) func(c *gin.Context, err error) {
	return func(c *gin.Context, err error) {
		c.AbortWithError(statuscode, err)
	}
}

type paramErr struct {
	message string
}

func (e paramErr) Error() string {
	return e.message
}
