package router

import (
	"FileReader/controllers"
	"FileReader/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	fmt.Println("hello peter")
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LogIn)
	r.POST("/file", middleware.MiddleWare(), controllers.FileReader)
	r.GET("/file", middleware.MiddleWare(), controllers.GetAll)
	r.PUT("/file", middleware.MiddleWare(), controllers.UpdateData)
	r.DELETE("/file/:id", middleware.MiddleWare(), controllers.DeleteData)
	r.GET("/joke", controllers.JokeHandler)
	r.POST("/refresh", controllers.Refresh)
	r.GET("/health", controllers.HealthHandler)
	r.GET("/readiness", controllers.ReadinessHandler)
	r.GET("/db-readiness", controllers.DBReadinessHandler)
}
