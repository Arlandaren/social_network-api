package router

import (
	"solution/api"
	// "solution/middlewares"
	"github.com/gin-gonic/gin"
)

func RouteAll(r * gin.Engine){
	api_router := r.Group("api")
	{
		api_router.GET("/ping", func(c *gin.Context){c.String(200, "pong")})
		api_router.GET("/countries", api.GetAllCountries)
		api_router.GET("/countries/:alpha2", api.GetCountryByid)
		// api_router.Use(middlewares.AuthValidation)
		auth := r.Group("auth")
		{
		auth.POST("/sign-in", api.Signin)
		auth.POST("/register",api.Register)
		}
	}
	
}