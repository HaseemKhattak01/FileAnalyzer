package main

import (
	"FileReader/router"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	// database.ConnectDatabase()
	router.Routes(g)
	g.Run(":1323")
}
