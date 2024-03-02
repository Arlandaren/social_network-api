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
		c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
		return
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if user.Username == "" || user.Username == "my" || !re.MatchString(user.Username){
		c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
		return
	}
	re = regexp.MustCompile(`\+[\d]+`)
	if !re.MatchString(user.PhoneNumber){
		c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
		return
	}
	existingUser, _ := models.GetUser(user.Username)
	if existingUser != nil && existingUser.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "user already exists"})
		return
	}
	if err := utils.CheckPassword(user.Password); err != nil {
		c.JSON(400, gin.H{"reason": err.Error()})
		return
	}
	hashedPassword, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "could not generate hash"})
		return
	}
	user.Password = hashedPassword

	profile, err := models.CreateUser(user.Username, user.Email, user.Password, user.CountryCode, user.PublicProfile, user.PhoneNumber, user.Image)
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
		c.JSON(409, gin.H{"reason": "Нарушена уникальность"})
		return
	}

	c.JSON(201, gin.H{"profile": profile})
}
func Signin(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
		return
	}
	existingUser, err := models.GetUser(user.Username)
	if err != nil {
		c.JSON(401, gin.H{"reason": "user not found"})
		return
	}

	if existingUser.ID == 0 {
		c.JSON(401, gin.H{"reason": "User not found"})
		return
	}
	if !utils.CompareHashPassword(user.Password, existingUser.Password) {
		c.JSON(400, gin.H{"reason": "Wrong password"})
		return
	}
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := models.Claims{
		User_id:        existingUser.ID,
		User_login: 	existingUser.Username,
		StandardClaims: jwt.StandardClaims{Subject: existingUser.Username, ExpiresAt: expirationTime.Unix()},
	}

	err = godotenv.Load()
	if err != nil {
		c.JSON(500, gin.H{"reason": "could not load environment variables"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		c.JSON(500, gin.H{"reason": "couldnt generate token"})
		return
	}
	c.SetCookie("token", signedString, int(expirationTime.Unix()), "/", strings.Split(os.Getenv("SERVER_ADDRESS"), ":")[0], false, true)
	c.JSON(200, gin.H{"token": signedString})
}
