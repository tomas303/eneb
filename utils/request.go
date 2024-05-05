package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryS(c *gin.Context, name string, defval string) string {
	v := c.Query(name)
	if v == "" {
		v = defval
	}
	return v
}

func QueryI(c *gin.Context, name string, defval int64) (int64, error) {
	v := QueryS(c, name, strconv.FormatInt(defval, 10))
	vi, err := strconv.ParseInt(v, 10, 64)
	return vi, err
}
