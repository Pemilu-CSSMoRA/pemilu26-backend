package routes

import (
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/controller"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/entity"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Election(route *gin.Engine, electionController controller.ElectionController, middleware middleware.Middleware) {
	routes := route.Group("/api/v1")
	{
		// Public
		routes.GET("/elections", electionController.GetAllElection)

		// Admin
		electionRoutes := routes.Group("/elections")
		electionRoutes.Use(middleware.Authenticate())
		{
			electionRoutes.POST("", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.CreateElection)
			electionRoutes.PUT("/:id", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.UpdateElection)
			electionRoutes.PUT("/:id/start", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.StartElection)
			electionRoutes.PUT("/:id/end", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.EndElection)
			electionRoutes.GET("/:id/voters", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.GetElectionVoters)
			electionRoutes.POST("/:id/voters", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.EnrollVoters)
			electionRoutes.DELETE("/:id/voters", middleware.OnlyAllow(string(entity.RoleAdmin)), electionController.UnenrollVoters)
		}
	}
}
