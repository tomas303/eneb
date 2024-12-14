package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type paramErr struct {
	message string
}

func (e paramErr) Error() string {
	return e.message
}

func ctxQParamStr(c *gin.Context, name string) *string {
	v, e := c.GetQuery(name)
	if e {
		return &v
	} else {
		return nil
	}
}

func ctxQParamInt(c *gin.Context, name string) *int {
	v, e := c.GetQuery(name)
	if e {
		x, err := strconv.Atoi(v)
		if err != nil {
			return nil
		} else {
			return &x
		}
	} else {
		return nil
	}
}

func ctxPParam(c *gin.Context, name string) *any {
	v, e := c.Get(name)
	if e {
		return &v
	} else {
		return nil
	}
}
