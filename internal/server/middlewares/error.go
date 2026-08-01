// Package middlewares contains the server middlewares
package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/nuttyshrimp/docker-dashboard/pkg/logger"
	"go.uber.org/zap"
)

const errMsg = "Server fout. Team Software is op de hoogte gebracht."

// nolint:gocognit // Keeping middlewares in one function
func ErrorHandler() fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
		}

		if code < 500 {
			if fiberErr != nil {
				return c.Status(code).SendString(fiberErr.Message)
			}
			return c.SendStatus(code)
		}

		log := zap.S()
		if logger.LocalLogger != nil {
			log = logger.LocalLogger.Sugar()
		}
		log.Errorf("%+v", err)

		return c.Status(500).SendString(errMsg)
	}
}
