package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nuttyshrimp/docker-dashboard/internal/server/service"
	"github.com/nuttyshrimp/docker-dashboard/pkg/ui"
	"github.com/nuttyshrimp/docker-dashboard/views"
)

type Views struct {
	router  fiber.Router
	service *service.Service
}

func NewViews(router fiber.Router, service *service.Service) *Views {
	api := &Views{
		router,
		service,
	}

	api.RegisterRoutes()

	return api
}

func (v *Views) RegisterRoutes() {
	v.router.Get("/", v.handleIndex)
}

func (v *Views) handleIndex(c fiber.Ctx) error {

	return ui.Render(c, views.List())
}
