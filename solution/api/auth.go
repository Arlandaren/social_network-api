package api

import (
	"fmt"
	"net/http"
	"os"
	"solution/models"
	"solution/utils"
	"time"
    "strings"
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
    fmt.Println(user.ID)

    existingUser, err := models.GetUser(user.Username)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if existingUser.ID != 0 {
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

    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "user created"})
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
        c.JSON(400, gin.H{"status": "error", "message": "couldnt generate token"})
        return
    }
    c.SetCookie("token", signedString,int(expirationTime.Unix()),"/",strings.Split(os.Getenv("SERVER_ADDRESS"), ":")[0],false,true)
    c.JSON(200,gin.H{"token":signedString})
}
// func Login(c *gin.Context) {
// 	var user models.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(400, gin.H{"status": "error", "message": err.Error()})
// 		return
// 	}
// 	var existingUser models.User
// 	models.DB.Where("name = ?", user.Name).First(&existingUser)
// 	if existingUser.ID == 0 {
// 		c.JSON(400, gin.H{"status": "error", "message": "user not exist"})
// 		return
// 	}
// 	errhash := utils.CompareHashPassword(user.Password, existingUser.Password)
// 	if !errhash {
// 		c.JSON(400, gin.H{"status": "error", "message": "invalid password"})
// 		return
// 	}
// 	expirationTime := time.Now().Add(720 * time.Hour)
// 	claims := &models.Claims{
// 		UserID: existingUser.ID,
// 		StandardClaims: jwt.StandardClaims{
// 			Subject:   existingUser.Name,
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}
// 	err := godotenv.Load()
// 	if err != nil {
// 		c.JSON(500, gin.H{"status": "error", "message": "could not load environment variables"})
// 		return
// 	}
// 	jwtkey := []byte(os.Getenv("JWT_KEY"))
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtkey)
// 	if err != nil {
// 		c.JSON(400, gin.H{"status": "error", "message": "couldnt generate token"})
// 		return
// 	}
// 	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
// 	c.JSON(200, gin.H{"status": "success", "message": "authentication success", "token": tokenString})
// }

// func ResetPassword(c *gin.Context) {

// 	var user models.User

// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var existingUser models.User

// 	models.DB.Where("email = ?", user.Name).First(&existingUser)

// 	if existingUser.ID == 0 {
// 		c.JSON(400, gin.H{"error": "user does not exist"})
// 		return
// 	}

// 	var errHash error
// 	user.Password, errHash = utils.GenerateHashPassword(user.Password)

// 	if errHash != nil {
// 		c.JSON(500, gin.H{"error": "could not generate password hash"})
// 		return
// 	}

// 	models.DB.Model(&existingUser).Update("password", user.Password)

// 	c.JSON(200, gin.H{"success": "password updated"})
// }