package main

import (
	"fmt"
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
	err = models.InitDB(os.Getenv("POSTGRES_CONN"))
	if err != nil{
		fmt.Println(err)
	}
	// Logger.Info("DB inititalized succesfully")
	r := gin.Default()
	router.RouteAll(r)

	r.Run(os.Getenv("SERVER_ADDRESS"))
}


 