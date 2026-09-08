package routes

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/controller"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, middleware middleware.Middleware) {
	routes := route.Group("/api/v1")
	{
		// Public
		routes.POST("/auth/register", userController.Register)
		routes.POST("/auth/login", userController.Login)

		// Authenticated user
		userRoutes := routes.Group("/users")
		userRoutes.Use(middleware.Authenticate())
		{
			// user
			userRoutes.GET("/me", userController.GetMe)
			userRoutes.PUT("/me", userController.UpdateMe)

			// admin
			userRoutes.PUT("/:id", middleware.OnlyAllow(string(entity.RoleAdmin)), userController.UpdateAdmin)
			//routes.DELETE("")
			userRoutes.GET("", middleware.OnlyAllow(string(entity.RoleAdmin)), userController.GetAllUser)

		}

		//routes.GET("/verification-status", userController.CheckVerificationStatus)
		//routes.GET("/me", middleware.Authenticate(), userController.Me)
		//routes.POST("/reset", userController.ResetPassword)
		//routes.POST("/forget", userController.ForgetPassword)
	}
}
