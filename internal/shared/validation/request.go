package validation

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, dst interface{}) error {
	return c.ShouldBindJSON(dst)
}

func ParsePositiveInt64Param(c *gin.Context, key string) (int64, error) {
	value, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, errors.New("value must be greater than zero")
	}

	return value, nil
}
