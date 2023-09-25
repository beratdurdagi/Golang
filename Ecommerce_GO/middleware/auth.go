package middleware

import (
	"net/http"

	jwttoken "github.com/Karalakrepp/Golang/Ecommerce_GO/jwt_token"
	"github.com/gin-gonic/gin"
)

func JWT_Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get Token from Request Header
		token := c.Request.Header.Get("token")

		//check if token doesn't exist
		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		//Get Claims
		claims, msg := jwttoken.ValidateToken(token)

		//Check Error or someting
		if msg == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't find claims"})
			c.Abort()
			return
		}

		//Set this
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		//Use middileware c.next()
		c.Next()

	}
}
