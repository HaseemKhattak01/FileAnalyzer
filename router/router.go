package router

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.POST("/signup", SignUp)
	r.POST("/login", LogIn)
	r.POST("/file", FileReader)
	r.GET("/file", GetAll)
	r.PUT("/file", UpdateData)
	r.DELETE("/file/:id", DeleteData)
}
