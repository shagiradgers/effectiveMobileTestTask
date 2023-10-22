package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetFromQueryInt - returns value, is exists, error
func GetFromQueryInt(ctx *gin.Context, key string) (int, bool, error) {
	v, ok := ctx.GetQuery(key)
	if !ok {
		return 0, false, nil
	}

	num, err := strconv.Atoi(v)
	if err != nil {
		return 0, true, err
	}
	return num, true, nil
}

func GetFromQueryUint(ctx *gin.Context, key string) (uint, bool, error) {
	v, ok := ctx.GetQuery(key)
	if !ok {
		return 0, false, nil
	}

	num, err := strconv.ParseUint(v, 10, strconv.IntSize)
	if err != nil {
		return 0, false, err
	}
	return uint(num), true, nil
}
