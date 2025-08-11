package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gerps2/desafio-cloud-run/shared/config"
	"github.com/gerps2/desafio-cloud-run/shared/logger"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	server *http.Server
	config *config.Config
	logger logger.Logger
}

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			RespondWithTimeout(c, "", []string{"Request exceeded the configured timeout"})
			c.Abort()
		}
	})
}

func ErrorHandlerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
				
				if !c.Writer.Written() {
					causes := []string{"An unexpected error occurred in the application"}
					RespondWithInternalError(c, "Internal server error", causes)
				}
				
				c.Abort()
			}
		}()
		
		c.Next()
		
		if len(c.Errors) > 0 {
			if !c.Writer.Written() {
				lastError := c.Errors.Last()
				logger.Error("Request error: %v", lastError.Error())
				
				causes := []string{lastError.Error()}
				RespondWithInternalError(c, "An error occurred while processing the request", causes)
			}
		}
	})
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	
	router.Use(ErrorHandlerMiddleware(log))
	
	timeout := time.Duration(cfg.App.RequestTimeoutSec) * time.Second
	router.Use(TimeoutMiddleware(timeout))

	return &Server{
		config: cfg,
		logger: log,
		router: router,
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info("Starting server on %s", addr)
	
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	return s.server.Shutdown(ctx)
}
