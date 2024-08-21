package main

import (
	"FileReader/database"
	"FileReader/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hey there..")
	g := gin.Default()
	database.ConnectDatabase()
	router.Routes(g)
	g.Run(":8080")
}
