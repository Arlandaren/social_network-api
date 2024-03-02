package router

import (
	"github.com/gin-gonic/gin"
	"solution/pkg/api"
	"solution/pkg/middlewares"
)

func RouteAll(r *gin.Engine) {
	api_router := r.Group("api")
	{
		api_router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
		api_router.GET("/countries", api.GetAllCountries)
		api_router.GET("/countries/:alpha2", api.GetCountryByid)
		auth := api_router.Group("auth")
		{
			auth.POST("/sign-in", api.Signin)
			auth.POST("/register", api.Register)
		}
		profile := api_router.Group("me")
		{
			profile.Use(middlewares.AuthValidation)
			profile.GET("/profile", api.Me)
			profile.PATCH("/profile", api.EditMe)
			profile.POST("/updatePassword", api.UpdatePassword)
		}
		friends := api_router.Group("friends")
		{
			friends.Use(middlewares.AuthValidation)
			friends.POST("/add", api.AddFriend)
			friends.POST("/remove", api.RemoveFriend)
			friends.GET("", api.GetFriendsList)
		}
		posts := api_router.Group("posts")
		{
			posts.Use(middlewares.AuthValidation)
			posts.POST("/new", api.NewPost)
			posts.GET("/:postId", api.GetPostById)
			feed := posts.Group("feed")
			{
				feed.GET("/:login", api.GetFeedByLogin)
				feed.GET("/my", api.GetMyFeed)
			}
		}
		api_router.Use(middlewares.AuthValidation)
		api_router.GET("/profiles/:login", api.Profiles)
	}

}
