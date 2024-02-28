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
	config := models.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBname:   os.Getenv("POSTGRES_DATABASE"),
	}
	err = models.InitDB(config)
	if err != nil {
		Logger.Fatal(err, " CONN- ", os.Getenv("POSTGRES_CONN"))
		// fmt.Println(err)
	}
	// Logger.Info("DB inititalized succesfully")
	r := gin.Default()
	router.RouteAll(r)

	r.Run(os.Getenv("SERVER_ADDRESS"))
}
