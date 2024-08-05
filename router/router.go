package router

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.POST("/home1", Signup)
	r.POST("/home2", Login)
	r.POST("/home", FileReader)
	r.GET("/home", Getall)
	r.PUT("/home", Updatedata)
	r.DELETE("/home/:id", Deletedata)
}
