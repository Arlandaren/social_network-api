package api

import (
	"solution/pkg/models"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func NewPost(c *gin.Context) {
	var postRequest models.PostRequest
	if err := c.ShouldBindJSON(&postRequest); err != nil {
		c.JSON(400, gin.H{"reason": err.Error()})
		return
	}
	if utf8.RuneCountInString(postRequest.Content) > 1000 {
		c.JSON(400, gin.H{"reason": "неверный формат"})
		return
	}
	for _, v := range postRequest.Tags {
		if utf8.RuneCountInString(v) > 20 {
			c.JSON(400, gin.H{"reason": "неверный формат"})
			return
		}
	}
	login, exists := c.Get("user_login")
	if !exists {
		c.JSON(401, gin.H{"reason": "Unauthorized"})
		return
	}
	postRequest.Author = login.(string)
	id, err := models.CreatePost(&postRequest)
	if err != nil {
		c.JSON(400, gin.H{"reason": err.Error()})
		return
	}
	post, err := models.GetPostById(id, postRequest.Author)
	if err != nil {
		c.JSON(404, gin.H{"reason": err.Error()})
		return
	}

	c.JSON(200, post)
}
func GetPostById(c *gin.Context) {
	login, exists := c.Get("user_login")
	if !exists {
		c.JSON(401, gin.H{"reason": "Unauthorized"})
		return
	}
	// postid,_ := strconv.ParseInt(c.Param("postId"), 10, 64)
	postId := c.Param("postId")
	post, err := models.GetPostById(postId, login.(string))
	if err != nil {
		c.JSON(404, gin.H{"reason": "пост не найден или к нему нет доступа"})
		return
	}
	c.JSON(200, post)
}
func GetMyFeed(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}

	if limit < 0 || limit > 50 || offset < 0 {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}
	login, exists := c.Get("user_login")
	if !exists {
		c.JSON(401, gin.H{"reason": "Unauthorized"})
		return
	}
	posts, err := models.GetMyFeedList(login.(string), offset, limit)
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
	if len(posts) == 0 {
		posts = make([]models.Post, 0)
	}
	c.JSON(200, posts)
}
func GetFeedByLogin(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}

	if limit < 0 || limit > 50 || offset < 0 {
		c.JSON(400, gin.H{"reason": "несоответствие формату"})
		return
	}

	targetLogin := c.Param("login")
	userLogin, exists := c.Get("user_login")
	if !exists {
		c.JSON(401, gin.H{"reason": "Unauthorized"})
		return
	}
	posts, err := models.GetFeedById(userLogin.(string), targetLogin, offset, limit)
	if err != nil {
		c.JSON(404, gin.H{"reason": "пост не найден или к нему нет доступа"})
		return
	}
	if len(posts) == 0 {
		posts = make([]models.Post, 0)
	}
	c.JSON(200, posts)
}
func LikePost(c *gin.Context) {
	userLogin, exists := c.Get("user_login")
	if !exists {
		c.JSON(401, gin.H{"reason": "Unauthorized"})
		return
	}
	post_id := c.Param("postId")
	if err := models.Like(userLogin.(string), post_id); err != nil {
		c.JSON(404, gin.H{"reason": err.Error()})
		return
	}
	post, err := models.GetPostById(post_id,userLogin.(string))
	if err != nil{
		c.JSON(404, gin.H{"reason": err.Error()})
		return
	}
	c.JSON(200,post)
}
func DislikePost(c *gin.Context) {
	userLogin, exists := c.Get("user_login")
	if !exists {
		c.JSON(401, gin.H{"reason": "Unauthorized"})
		return
	}
	post_id := c.Param("postId")
	if err := models.Dislike(userLogin.(string), post_id); err != nil {
		c.JSON(404, gin.H{"reason1": err.Error()})
		return
	}
	post, err := models.GetPostById(post_id,userLogin.(string))
	if err != nil{
		c.JSON(404, gin.H{"reason2": err.Error()})
		return
	}
	c.JSON(200,post)
}
