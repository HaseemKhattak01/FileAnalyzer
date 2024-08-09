package middleware

import (
	"FileReader/Jwt"
	"FileReader/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MiddleWare() gin.HandlerFunc {
	return func(g *gin.Context) {
		apitoken := g.Request.Header.Get("Authorization")
		err := Jwt.Validity(apitoken)
		if err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		g.Next()
	}
}
