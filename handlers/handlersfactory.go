package handlers

import (
	"eneb/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = func(*gin.Context)

func MakeHandlerGetOne[T any](param ParamGetterFunc, cmd data.DataCmdSelectOneFunc[T]) HandlerFunc {
	return func(c *gin.Context) {
		id := param(c)
		value, err := cmd(id)
		if err != nil {
			c.Set("error", err)
			return
		}
		c.IndentedJSON(http.StatusOK, value)
	}
}

func MakeHandlerGetMany[T any](params []ParamGetterFunc, cmd data.DataCmdSelectManyFunc[T], reverse bool) HandlerFunc {
	return func(c *gin.Context) {
		args := make([]any, len(params))
		for i := 0; i < len(params); i++ {
			args[i] = params[i](c)
		}
		value, err := cmd(args)
		if err != nil {
			c.Set("error", err)
			return
		}
		if reverse {
			for i := 0; i < len(value)/2; i++ {
				value[i], value[len(value)-1-i] = value[len(value)-1-i], value[i]
			}
		}
		c.IndentedJSON(http.StatusOK, value)
	}
}

func MakeHandlerPostOne[T any](cmd data.DataCmdSaveOneFunc[T]) HandlerFunc {
	return func(c *gin.Context) {
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
