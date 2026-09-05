package routes

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/controller"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, middleware middleware.Middleware) {
	routes := route.Group("/api/auth")
	{
		// User
		routes.POST("", userController.Register)
		routes.POST("/login", userController.Login)
		//routes.PUT("/update", middleware.Authenticate(), userController.Update)

		//routes.POST("/sendmail", userController.SendVerificationEmail)
		//routes.GET("/verify", userController.VerifyEmail)

		//routes.GET("", middleware.Authenticate(), middleware.OnlyAllow(string(entity.RoleAdmin)), userController.GetAllUser)
		//routes.GET("/verification-status", userController.CheckVerificationStatus)
		//routes.GET("/me", middleware.Authenticate(), userController.Me)
		//routes.POST("/reset", userController.ResetPassword)
		//routes.POST("/forget", userController.ForgetPassword)
	}
}
