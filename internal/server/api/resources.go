package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nuttyshrimp/docker-dashboard/internal/server/service"
)

type Resources struct {
	router  fiber.Router
	service *service.Service
}

func NewResources(router fiber.Router, service *service.Service) *Resources {
	api := &Resources{
		router:  router.Group("/resources"),
		service: service,
	}
	return api
}
