package router

import (
	"solution/api"
	"solution/middlewares"
	"github.com/gin-gonic/gin"
)

func RouteAll(r * gin.Engine){
	api_router := r.Group("api")
	{
		api_router.GET("/ping", func(c *gin.Context){c.JSON(200, gin.H{"status":"ok"})})
		api_router.GET("/countries", api.GetAllCountries)
		api_router.GET("/countries/:alpha2", api.GetCountryByid)
		api_router.Use(middlewares.AuthValidation)
	}
	auth := r.Group("auth")
	{
		auth.GET("/sing-in")
		auth.POST("/register")
	}
}