package controllers

import "github.com/gofiber/fiber/v2"

func (s *Service) getAllSystemRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupHealthRoute,
	}
}

func (s *Service) setupHealthRoute(routeGroup fiber.Router) {
	routeGroup.Get("/health", s.healthCheck)
}

func (s *Service) healthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ok",
		"service": "umineko-quote-service",
	})
}
