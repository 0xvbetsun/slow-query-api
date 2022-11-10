package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Server represents API server of application
type Server struct {
	app *fiber.App
}

// Run creates configuration for the server and starts it
func (s *Server) Run(addr string, h *Handler) error {
	app := fiber.New(fiber.Config{})
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/queries", h.GetAll)

	return s.app.Listen(addr)
}

// Shutdown stops the Server
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
