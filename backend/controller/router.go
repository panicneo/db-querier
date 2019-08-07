package controller

import (
	"db-querier/controller/auth"
	"db-querier/controller/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

// register Handlers
func Register(engine *gin.Engine) {
	engine.GET("", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/static")
	})
	engine.Static("/static", "./dist")
	api := engine.Group("api")
	api.Use(func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) != 0 {
			ctx.AbortWithStatusJSON(400, ctx.Errors.JSON())
		}
	})
	authG := api.Group("auth")
	{
		authG.POST("/login", auth.Login)
	}
	queryG := api.Group("query")
	{
		queryG.GET("", query.List)
		queryG.GET(":name", query.Query)
	}
}
