package routers

import (
	"go-project-template/app/api"
	"go-project-template/app/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	r = ApiV0(r)
	return r
}

func ApiV0(r *gin.Engine) *gin.Engine {
	// r.Static("/upload", "public/upload")

	v0 := r.Group("/api").Use(middleware.TestUser)
	{
		v0.GET("/hello", api.HelloController{}.Hello)

	}

	return r
}
