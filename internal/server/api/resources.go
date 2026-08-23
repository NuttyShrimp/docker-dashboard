package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nuttyshrimp/docker-dashboard/internal/server/service"
)

type Resources struct {
	router    fiber.Router
	service   *service.Service
	subdomain *service.Subdomain
}

func NewResources(router fiber.Router, service *service.Service) *Resources {
	api := &Resources{
		router:    router.Group("/resources"),
		service:   service,
		subdomain: service.NewSubdomain(),
	}
	api.Routes()
	return api
}

func (r *Resources) Routes() {
	r.router.Get("/", r.handleGetAll)
}

func (r *Resources) handleGetAll(c fiber.Ctx) error {
	return c.JSON(r.subdomain.GetAll())
}
