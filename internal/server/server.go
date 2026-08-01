// Package server starts the server
package server

import (
	"fmt"

	zapfiber "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/nuttyshrimp/docker-dashboard/internal/database/repository"
	routers "github.com/nuttyshrimp/docker-dashboard/internal/server/api"
	middlewares "github.com/nuttyshrimp/docker-dashboard/internal/server/middlewares"
	"github.com/nuttyshrimp/docker-dashboard/internal/server/service"
	"github.com/nuttyshrimp/docker-dashboard/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	*fiber.App
	Addr string
}

func New() *Server {
	repo := repository.New()
	service := service.New(repo)

	// Construct app
	app := fiber.New(fiber.Config{
		BodyLimit:         20 * 1024 * 1024,
		ReadBufferSize:    8096,
		StreamRequestBody: true,
		ErrorHandler:      middlewares.ErrorHandler(),
	})

	app.Use(recover.New())
	app.Use(zapfiber.New(zapfiber.Config{
		Logger: zap.L(),
	}))
	if config.IsDev() {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Access-Control-Allow-Origin"},
			AllowCredentials: true,
		}))
	}

	// Register routes
	api := app.Group("/api")
	routers.NewResources(api, service)

	// Fallback
	app.All("/api*", func(c fiber.Ctx) error {
		return c.SendStatus(404)
	})

	port := config.GetDefaultInt("server.port", 8000)
	host := config.GetDefaultString("server.host", "0.0.0.0")

	srv := &Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
		App:  app,
	}

	return srv
}
