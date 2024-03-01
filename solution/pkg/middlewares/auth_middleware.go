package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthValidation(c *gin.Context){
	header := c.GetHeader("Authorization")
	if header == ""{
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}
	tokenString := strings.Replace(header, "Bearer ", "", 1)
	err := godotenv.Load()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "could not load environment variables"})
			return
		}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){return []byte(os.Getenv("JWT_KEY")),nil})
	if err != nil{
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"error with parsing token", "message":err.Error()})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
		return
	}
	userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "message": err.Error()})
			return
		}
	c.Set("user_id", uint(userID))
}