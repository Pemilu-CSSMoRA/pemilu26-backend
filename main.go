package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Pemilu-CSSMoRA/pemilu26-backend/cmd"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/config"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/dto"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/logger"
	"github.com/Pemilu-CSSMoRA/pemilu26-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadEnv(); err != nil {
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

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Pemilu 2026 Backend API",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	port := config.GetEnv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	router.Run(":" + port)
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
