package api

import (
	"net/http"
	"os"
	"regexp"
	"solution/pkg/models"
	"solution/pkg/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	existingUser, _ := models.GetUser(user.Username)
	if existingUser != nil && existingUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "user already exists"})
		return
	}
	regex := regexp.MustCompile(`[a-z].*[A-Z]|[A-Z].*[a-z]|[0-9].*[a-zA-Z]|[a-zA-Z].*[0-9]`)
	if !regex.MatchString(user.Password) {
		c.JSON(400, "Регистрационные данные не соответствуют ожидаемому формату и требованиям.")
		return
	}
    regex = regexp.MustCompile(`[A-Z]`)
    if !regex.MatchString(user.Password) {
        c.JSON(400, "Пароль должен содержать хотя бы одну букву в верхнем регистре")
        return
    }
	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "could not generate hash"})
		return
	}
	user.Password = hashedPassword

	profile, err := models.CreateUser(user.Username, user.Email, user.Password, user.CountryCode, user.PublicProfile, user.PhoneNumber, user.Image)
	if err != nil {
        if err.Error() == "неверный формат" {
            c.JSON(400, "Регистрационные данные не соответствуют ожидаемому формату и требованиям.")
            return
        }
		if err.Error() == "pq: значение не умещается в тип character varying(200)" {
			c.JSON(400, gin.H{"error": "Длинна ссылки на аватар больше 200", "message": "could not create user"})
			return
		}
		if err.Error() == "pq: значение не умещается в тип character varying(2)" || err.Error() == "pq: INSERT или UPDATE в таблице \"users\" нарушает ограничение внешнего ключа \"fk_country\"" {
			c.JSON(400, gin.H{"error": "Не найдено такого кода страны", "message": "could not create user"})
			return
		}
        if err.Error() == "pq: значение не умещается в тип character varying(20)" || err.Error() == "pq: INSERT или UPDATE в таблице \"users\" нарушает ограничение внешнего ключа \"fk_country\"" {
			c.JSON(400, gin.H{"error": "Не найдено такого кода страны", "message": "could not create user"})
			return
		}
		c.JSON(409, gin.H{"error": err.Error(), "message": "could not create user"})
		return
	}

	c.JSON(201, gin.H{"profile": profile})
}
func Signin(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	existingUser, err := models.GetUser(user.Username)
	if err != nil {
		c.JSON(401, gin.H{"error": "user not found", "message": err.Error()})
		return
	}

	if existingUser.ID == 0 {
		c.JSON(401, gin.H{"error": "User not found", "message": err.Error()})
		return
	}
	if !utils.CompareHashPassword(user.Password, existingUser.Password) {
		c.JSON(400, gin.H{"error": "Wrong password"})
		return
	}
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := models.Claims{
		User_id:        existingUser.ID,
		StandardClaims: jwt.StandardClaims{Subject: existingUser.Username, ExpiresAt: expirationTime.Unix()},
	}

	err = godotenv.Load()
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "could not load environment variables"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "message": "couldnt generate token"})
		return
	}
	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", strings.Split(os.Getenv("SERVER_ADDRESS"), ":")[0], false, true)
	c.JSON(200, gin.H{"token": signedString})
}
