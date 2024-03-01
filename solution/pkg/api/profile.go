package api

import (
	"net/http"
	"solution/pkg/models"
	"solution/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Unauthorized"})
		return 
	}
	profile,err := models.GetMyProfile(userId.(uint))
	if err != nil{
		c.JSON(500,gin.H{"reason": "не удалось получить профиль"})
		return
	}
	c.JSON(200,profile)
}
func UpdatePassword(c *gin.Context){
	userId,exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Unauthorized"})
		return 
	}
	var updatePasswordForm *models.UpdatePasswordForm
	if err := c.ShouldBindJSON(&updatePasswordForm); err != nil{
		c.JSON(400,gin.H{"reason":"Новый пароль не соответствует требованиям безопасности."})
		return
	}
	if err := utils.CheckPassword(updatePasswordForm.NewPasword);err!=nil{
		c.JSON(400,gin.H{"reason":"Новый пароль не соответствует требованиям безопасности."})
		return
	}
	if err:=models.UpdatePassword(updatePasswordForm, userId.(uint)); err != nil{
		c.JSON(403, gin.H{"reason":err.Error()})
		return
	}
	header := c.GetHeader("Authorization")
	if header == ""{
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Auth header is missing"})
		return
	}
	tokenString := strings.Replace(header, "Bearer ", "", 1)
	
	models.DeactivateToken(tokenString)

	c.JSON(200,gin.H{"status":"ok"})
}
func EditMe(c *gin.Context){
	var editParameters models.EditParameters
	if err:=c.ShouldBindJSON(&editParameters); err !=nil{
		c.JSON(400,gin.H{"reason":"несоответствие формату"})
		return
	}
	userId, exists := c.Get("user_id")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"reason":"Unauthorized"})
		return 
	}

	err := models.UpdateProfile(userId.(uint), &editParameters)
	if err != nil {
		if err.Error() == "неверный формат" {
			c.JSON(400, gin.H{"reason": "Регистрационные данные не соответствуют ожидаемому формату и требованиям."})
			return
		}
		if err.Error() == "pq: значение не умещается в тип character varying(200)" {
			c.JSON(400, gin.H{"reason": "Длинна ссылки на аватар больше 200"})
			return
		}
		if err.Error() == "pq: значение не умещается в тип character varying(2)" || err.Error() == "pq: INSERT или UPDATE в таблице \"users\" нарушает ограничение внешнего ключа \"fk_country\"" {
			c.JSON(400, gin.H{"reason": "Не найдено такого кода страны"})
			return
		}
		if err.Error() == "pq: значение не умещается в тип character varying(20)" || err.Error() == "pq: INSERT или UPDATE в таблице \"users\" нарушает ограничение внешнего ключа \"fk_country\"" {
			c.JSON(400, gin.H{"reason": "Номер телефона слишком длинный"})
			return
		}
		c.JSON(409, gin.H{"reason": "нарушена уникальность"})
		return
	}
	profile,err := models.GetMyProfile(userId.(uint))
	if err != nil{
		c.JSON(500,gin.H{"reason":"не удалось получить профиль"})
		return
	}
	c.JSON(200,profile)
}
func Profiles(c *gin.Context){
	login := c.Param("login")
	profile, err := models.GetProfile(login)
	if err != nil{
		if err.Error() == "sql: no rows in result set"{
			c.JSON(403, gin.H{"reason":"профиль не публичный"})
			return
		}
		c.JSON(400, gin.H{"reason":err.Error()})
		return
	}
	c.JSON(200,profile)
}

  