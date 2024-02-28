package main

import (
	"os"
	"solution/models"
	"solution/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	Logger := NewLogger()

	gin.DefaultWriter = Logger.Writer()

	err := godotenv.Load()

	if err != nil {
		Logger.Fatal("Couldnt load env variables")
	}
	dsn := os.Getenv("POSTGRES_CONN")+"?sslmode=disable"
	err = models.InitDB(dsn)
	if err != nil {
		Logger.Fatal(err)
	}
	r := gin.Default()
	router.RouteAll(r)

	r.Run(os.Getenv("SERVER_ADDRESS"))
}
