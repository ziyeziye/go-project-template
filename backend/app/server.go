package app

import (
	"fmt"
	"go-project-template/app/routers"
	"go-project-template/config"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunServe() {
	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// 捕获 panic 错误
	app.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  http.StatusText(http.StatusInternalServerError),
		})
	}))

	app = routers.InitRouter(app)

	port := config.Get().Port
	if gin.Mode() != gin.ReleaseMode {
		// 打印 API 路由
		log.Println("Http API List:")

		for _, r := range app.Routes() {
			if r.Method != "OPTIONS" {
				fmt.Printf("[%s] http://127.0.0.1:%d%s\n", r.Method, port, r.Path)
			}
		}
	}

	log.Fatal(app.Run(fmt.Sprintf(":%d", port)))
}
