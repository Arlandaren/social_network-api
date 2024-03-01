package router

import (
	"solution/pkg/api"
	"solution/pkg/middlewares"
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
		auth := api_router.Group("auth")
		{
		auth.POST("/sign-in", api.Signin)
		auth.POST("/register",api.Register)
		}
		profile := api_router.Group("me")
		{
			profile.Use(middlewares.AuthValidation)
			profile.GET("/profile",api.Me)
		}
	}
	
}