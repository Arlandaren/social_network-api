package api

import (
	"net/http"
	"solution/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddFriend(c *gin.Context){
	login, exists := c.Get("user_login")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Unauthorized"})
		return 
	}
	type friendRequest struct {
		Login string `json:"login"`
	}	
	var friendLogin friendRequest
    if err := c.BindJSON(&friendLogin); err != nil {
        c.JSON(400, gin.H{"reason": err.Error()})
        return
    }
	if friendLogin.Login == login{
		c.JSON(200,gin.H{"status":"ok"})
		return
	}
	if err := models.AddFriend(friendLogin.Login, login.(string)); err!=nil{
		if err.Error() == "pq: повторяющееся значение ключа нарушает ограничение уникальности \"friendships_friend_login_key\""{
			c.JSON(200,gin.H{"status":"ok"})
			return
		}
		c.JSON(404,gin.H{"reason":"User not found"})
		return
	}
	c.JSON(200,gin.H{"status":"ok"})
}
func RemoveFriend(c *gin.Context){
	login, exists := c.Get("user_login")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Unauthorized"})
		return 
	}
	var friendLogin models.FriendRequest
    if err := c.BindJSON(&friendLogin); err != nil {
        c.JSON(400, gin.H{"reason": "User nor found"})
        return
    }
	if err := models.RemoveFriend(friendLogin.Login, login.(string)); err!=nil{
		c.JSON(404,gin.H{"reason":"User not found"})
		return
	}
	c.JSON(200,gin.H{"status":"ok"})
}
func GetFriendsList(c *gin.Context){
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{"reason":"несоответствие формату"})
		return
	}

	if limit < 0 || limit > 50 || offset < 0 {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}
	login, exists := c.Get("user_login")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Unauthorized"})
		return 
	}
	friends,err:=models.GetFriendsList(login.(string),offset,limit)
	if err!=nil{
		c.JSON(200, friends)
		return
	}
	if len(friends) == 0{
		friends = make([]models.Friend, 0)
	}
	c.JSON(200, friends)
}