package ui

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/nuttyshrimp/docker-dashboard/views/layouts"
)

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return layouts.BaseLayout(component).Render(c.Context(), c.Response().BodyWriter())
}
