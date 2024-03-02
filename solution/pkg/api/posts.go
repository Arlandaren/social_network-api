package api

import (
	"solution/pkg/models"

	"github.com/gin-gonic/gin"
)

func NewPost(c *gin.Context){
	var postRequest models.PostRequest
	if err :=  c.ShouldBindJSON(&postRequest); err !=nil{
		c.JSON(400,gin.H{"reason":err.Error()})
		return
	}
	login, exists := c.Get("user_login")
	if !exists{
		c.JSON(401,gin.H{"reason":"Unauthorized"})
		return 
	}
	postRequest.Author = login.(string)
	id,err := models.CreatePost(&postRequest)
	if err != nil{
		c.JSON(400,gin.H{"reason":err.Error()})
		return
	}
	post,err := models.GetPostById(id)
	if err != nil{
		c.JSON(404,gin.H{"reason":err.Error()})
		return
	}

	c.JSON(200,post)
}