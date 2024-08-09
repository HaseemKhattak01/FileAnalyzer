package router

import (
	"FileReader/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.POST("/signup", SignUp)
	r.POST("/login", LogIn)
	r.POST("/file", middleware.MiddleWare(), FileReader)
	r.GET("/file", middleware.MiddleWare(), GetAll)
	r.PUT("/file", middleware.MiddleWare(), UpdateData)
	r.DELETE("/file/:id", middleware.MiddleWare(), DeleteData)
}
