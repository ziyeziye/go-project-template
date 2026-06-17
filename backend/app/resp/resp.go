package resp

import (
	"go-project-template/app/errno"
	"go-project-template/common/logx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type H map[string]interface{}

func JsonOk(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func JsonErr(ctx *gin.Context, err error, logErr ...error) {
	if len(logErr) > 0 {
		if logErr[0] != nil {
			logx.WriteError(logErr[0])
		}
	}

	newErr := errno.ErrOf(err)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": newErr.Code,
		"msg":  newErr.Msg,
		"data": nil,
	})
}

func JsonErrWithMsg(ctx *gin.Context, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  msg,
		"data": nil,
	})
}

func JsonErrWithData(ctx *gin.Context, err error, data interface{}) {
	newErr := errno.ErrOf(err)
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": newErr.Code,
		"msg":  newErr.Msg,
		"data": data,
	})
}
