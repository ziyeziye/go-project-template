package api

import (
	"errors"
	"go-project-template/app/resp"
	"go-project-template/service"

	"github.com/gin-gonic/gin"
)

type HelloController struct{ BaseController }

func (s HelloController) Hello(ctx *gin.Context) {
	test := ctx.Query("test")
	if test == "error" {
		err := service.UserSvr.TestAbortError()
		if err != nil {
			resp.JsonErr(ctx, err, errors.New("test error"))
			return
		}
		return
	}

	resp.JsonOk(ctx, resp.H{
		"hello": ctx.MustGet("user_name"),
	})
}
