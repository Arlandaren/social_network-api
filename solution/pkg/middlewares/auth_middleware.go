package middlewares

import (
	"net/http"
	"os"
	"solution/pkg/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthValidation(c *gin.Context){
	header := c.GetHeader("Authorization")
	if header == ""{
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "Authorization header is missing"})
		return
	}
	tokenString := strings.Replace(header, "Bearer ", "", 1)
	count,err := models.CheckBlackList(tokenString)
	if err != nil{
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"reason": "error with db"})
	}
	if count > 0{
		c.AbortWithStatusJSON(401, gin.H{"reason": "Token invalid"})
	}
	err = godotenv.Load()
		if err != nil {
			c.JSON(500, gin.H{"reason": "could not load environment variables"})
			return
		}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){return []byte(os.Getenv("JWT_KEY")),nil})
	if err != nil{
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason":"error with parsing token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "Invalid token"})
		return
	}
	userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"reason": "Invalid token"})
			return
		}

	c.Set("user_id", uint(userID))

}