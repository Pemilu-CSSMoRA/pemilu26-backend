package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/cmd"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/config"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/controller"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/middleware"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/repository"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/routes"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/service"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/logger"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		logger.Errorf("Warning: %v", err.Error())
		logger.Errorf("Continuing without .env file. Ensure all necessary environment variables are set.")
	}

	logger.Infof("Setting up database connection...")
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)
	logger.Infof("Database connection established")

	// Check if there are any command-line arguments
	if len(os.Args) > 1 {
		logger.Infof("Running commands...")
		cmd.Commands(db)
		for _, arg := range os.Args[1:] {
			if arg == "--migrate" || arg == "--seed" || arg == "--list-scripts" ||
				strings.HasPrefix(arg, "--script=") || arg == "--help" {
				return
			}
		}
	}

	var (
		// repository
		userRepository repository.UserRepository = repository.NewUserRepository(db)

		// service
		jwtService        service.JWTService    = service.NewJWTService(cfg.JWTSecret, cfg.JWTExpireHours)
		middlewareService middleware.Middleware = middleware.New(db, jwtService)

		userService service.UserService = service.NewUserService(userRepository, jwtService)

		// controller
		userController controller.UserController = controller.NewUserController(userService)
	)

	logger.Infof("Services initialized")
	logger.Infof("Setting up server...")
	server := gin.Default()
	server.Use(handlePanic())
	server.MaxMultipartMemory = 30 * 1024 * 1024
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
		})
	})

	server.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pemilu 2026 Backend API",
		})
	})

	server.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	//routes
	routes.User(server, userController, middlewareService)
	port := config.GetEnv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	server.Run(":" + port)
}

func handlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}
				fmt.Printf("\n[recovery] panic occurred: %v\n", err)
				stack := debug.Stack()
				fmt.Fprintln(os.Stderr, string(stack))

				response.BuildResponseFailed(dto.MESSAGE_FAILED_PANIC_OCCURED, err).
					SendWithAbort(ctx)
			}
		}()

		ctx.Next()
	}
}
