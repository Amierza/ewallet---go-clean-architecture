package routes

import (
	"github.com/Amierza/e-wallet/controller"
	"github.com/Amierza/e-wallet/middleware"
	"github.com/Amierza/e-wallet/service"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	routes := route.Group("api/user")
	{
		// User
		routes.POST("/register", userController.Register)
		routes.POST("/login", userController.Login)
		routes.POST("/topup", middleware.Authenticate(jwtService), userController.TopUp)
		routes.POST("/pay", middleware.Authenticate(jwtService), userController.Payment)
		routes.POST("/transfer", middleware.Authenticate(jwtService), userController.Transfer)
		routes.GET("/get-all-user", middleware.Authenticate(jwtService), userController.GetAllUser)
		routes.GET("/transactions", middleware.Authenticate(jwtService), userController.GetAllTransaction)
		routes.POST("/update-profile", middleware.Authenticate(jwtService), userController.UpdateProfile)
	}
}
