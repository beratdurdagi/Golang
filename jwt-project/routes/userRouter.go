package routes

import (
	controller "JWT_PROJECT/controllers"

	"github.com/gin-gonic/gin"

	middleware "JWT_PROJECT/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())

}
