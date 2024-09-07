package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func queryS(c *gin.Context, name string, defval string) string {
	v := c.Query(name)
	if v == "" {
		v = defval
	}
	return v
}
func queryI64(c *gin.Context, name string, defval int64) (int64, error) {
	v := queryS(c, name, strconv.FormatInt(defval, 10))
	val, err := strconv.ParseInt(v, 10, 64)
	return val, err
}

func queryI(c *gin.Context, name string, defval int) (int, error) {
	v := queryS(c, name, strconv.Itoa(defval))
	val, err := strconv.Atoi(v)
	return val, err
}

func qParamAsInt(name string, defval int, onerror func(c *gin.Context, err error)) func(c *gin.Context) int {
	return func(c *gin.Context) int {
		val, err := queryI(c, name, defval)
		if err != nil {
			onerror(c, err)
		}
		return val
	}
}

func qParamAsInt64(name string, defval int64, onerror func(c *gin.Context, err error)) func(c *gin.Context) int64 {
	return func(c *gin.Context) int64 {
		val, err := queryI64(c, name, defval)
		if err != nil {
			onerror(c, err)
		}
		return val
	}
}
