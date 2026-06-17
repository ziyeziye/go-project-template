package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	ReqUserName = "user_name"
)

func getUser(ctx *gin.Context) string {
	token := ctx.GetHeader(ReqUserName)
	if token == "" {
		token = ctx.Query(ReqUserName)
	}
	if token == "" {
		token = ctx.PostForm(ReqUserName)
	}
	return token
}

func TestUser(ctx *gin.Context) {
	user := getUser(ctx)
	if user == "" {
		user = "unknown"
		// _ = ctx.Error(errno.Unauthorized)
		// ctx.Abort()
		// return
	}

	ctx.Set("user_name", user)
	ctx.Next()
}
