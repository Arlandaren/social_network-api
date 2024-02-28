package api

import (
	"fmt"
	"net/http"
	"os"
	"solution/models"
	"solution/utils"
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

    hashedPassword, err := utils.GenerateHashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "could not generate hash"})
        return
    }
    user.Password = hashedPassword

    err = models.CreateUser(user.Username,user.Email,user.Password,user.Country,user.PublicProfile,user.PhoneNumber,user.Image)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "could not create user"})
        return
    }

    c.JSON(201, gin.H{"profile":user})
}
func Signin(c *gin.Context){
	var user *models.User
	if err := c.ShouldBindJSON(&user); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
        return
    }
    existingUser, err := models.GetUser(user.Username)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
        return
    }

    if existingUser.ID == 0{
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User not found"})
        return
    }
    if !utils.CompareHashPassword(user.Password, existingUser.Password){
        c.JSON(500,gin.H{"error": "Wrong password"})
        return
    }
    expirationTime := time.Now().Add(60* time.Minute)
    claims := models.Claims{
        User_id: existingUser.ID,
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
        fmt.Println("||||||||||||||||||||||||||||||||||||||||||||||||||      ",[]byte(os.Getenv("JWT_KEY")), "||||||||||||||||||||||||||",err.Error())
        c.JSON(500, gin.H{"status": "error", "message": "couldnt generate token"})
        return
    }
    c.SetCookie("token", signedString,int(expirationTime.Unix()),"/",strings.Split(os.Getenv("SERVER_ADDRESS"), ":")[0],false,true)
    c.JSON(200,gin.H{"token":signedString})
}