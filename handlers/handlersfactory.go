package handlers

import (
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = func(*gin.Context, []any)

func MakeHandlerGetOne[T any](cmd data.DataCmdSelectOneFunc[T]) HandlerFunc {
	return func(c *gin.Context, params []any) {
		id := params[0]
		value, err := cmd(id)
		if err != nil {
			c.Set("error", err)
			return
		}
		c.IndentedJSON(http.StatusOK, value)
	}
}

func MakeHandlerGetMany[T any](cmd data.DataCmdSelectManyFunc[T]) HandlerFunc {
	return func(c *gin.Context, params []any) {
		value, err := cmd(params)
		if err != nil {
			c.Set("error", err)
			return
		}
		c.IndentedJSON(http.StatusOK, value)
	}
}

func MakeHandlerPostOne[T any](cmd data.DataCmdSaveOneFunc[T]) HandlerFunc {
	return func(c *gin.Context, params []any) {
		var en T
		if err := c.BindJSON(&en); err != nil {
			c.Set("error", err)
			return
		}
		en, err := cmd(en)
		if err != nil {
			c.Set("error", err)
			return
		}
		c.Status(http.StatusCreated)
	}
}
