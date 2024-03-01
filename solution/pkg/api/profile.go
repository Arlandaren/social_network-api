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
	profile,err := models.GetMyProfile(userId.(uint))
	if err != nil{
		c.JSON(500,gin.H{"error": "error with Db","message":err.Error()})
		return
	}
	c.JSON(200,profile)
}
func EditMe(c *gin.Context){
	var editParameters models.EditParameters
	if err:=c.ShouldBindJSON(&editParameters); err !=nil{
		c.JSON(400,gin.H{"error":"несоответствие формату"})
		return
	}
	userId, exists := c.Get("user_id")
	if !exists{
		c.JSON(http.StatusUnauthorized,"Unauthorized")
		return 
	}

	err := models.UpdateProfile(userId.(uint), &editParameters)
	if err != nil{
		if err.Error() == "неверный формат" {
            c.JSON(400, "Регистрационные данные не соответствуют ожидаемому формату и требованиям.")
            return
        }
		if err.Error() == "pq: значение не умещается в тип character varying(200)" {
			c.JSON(400, gin.H{"error": "Длинна ссылки на аватар больше 200", "message": "could update user"})
			return
		}
		if err.Error() == "pq: значение не умещается в тип character varying(2)" || err.Error() == "pq: INSERT или UPDATE в таблице \"users\" нарушает ограничение внешнего ключа \"fk_country\"" {
			c.JSON(400, gin.H{"error": "Не найдено такого кода страны", "message": "could update user"})
			return
		}
        if err.Error() == "pq: значение не умещается в тип character varying(20)" || err.Error() == "pq: INSERT или UPDATE в таблице \"users\" нарушает ограничение внешнего ключа \"fk_country\"" {
			c.JSON(400, gin.H{"error": "Не найдено такого кода страны", "message": "could update user"})
			return
		}
		c.JSON(409, gin.H{"error": "нарушена уникальность", "message": "could update user"})
		return
	}
	profile,err := models.GetMyProfile(userId.(uint))
	if err != nil{
		c.JSON(500,gin.H{"error":err.Error()})
		return
	}
	c.JSON(200,profile)
}
func Profiles(c *gin.Context){
	login := c.Param("login")
	profile, err := models.GetProfile(login)
	if err != nil{
		if err.Error() == "sql: no rows in result set"{
			c.JSON(403, gin.H{"error":"профиль не публичный"})
			return
		}
		c.JSON(400, gin.H{"error":err.Error()})
		return
	}
	c.JSON(200,profile)
}
  