package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerps2/desafio-cloud-run/features/weather"
	httpServer "github.com/gerps2/desafio-cloud-run/shared/http"
	"github.com/gerps2/desafio-cloud-run/shared/logger"

	"github.com/gin-gonic/gin"
)

type App struct {
	server            *httpServer.Server
	weatherController *weather.WeatherController
	logger            logger.Logger
}

func NewApp(
	server *httpServer.Server,
	weatherController *weather.WeatherController,
	logger logger.Logger,
) *App {
	return &App{
		server:            server,
		weatherController: weatherController,
		logger:            logger,
	}
}

func (a *App) setupRoutes() {
	router := a.server.GetRouter()

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Register weather routes
	a.weatherController.RegisterRoutes(router)
}

func (a *App) Run() error {
	a.setupRoutes()

	// Start server in a goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Error("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return a.server.Shutdown(ctx)
}

func main() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
