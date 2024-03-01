package api

import (
	"net/http"
	"solution/pkg/models"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists{
		c.JSON(http.StatusUnauthorized,"Unauthorized")
		return 
	}
	// userId, err := utils.GetUserId(c)
	// if err != nil || userId == 0{
	// 	c.JSON(http.StatusUnauthorized,"Unauthorized")
	// }
	profile,err := models.GetProfile(userId.(uint))
	if err != nil{
		c.JSON(500,gin.H{"error": "error with Db","message":err.Error()})
		return
	}
	c.JSON(200,profile)
}
