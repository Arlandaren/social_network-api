package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) (uint, error) {
	id, ok := c.Get("user_id")
	if !ok {
		return 0, errors.New("user id not found,unauthorized???")
	}

	idUInt, ok := id.(uint)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idUInt, nil
}