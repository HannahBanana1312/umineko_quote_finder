package routes

import (
	"umineko_quote/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(service controllers.Service, app *fiber.App) {
	allRoutes := service.GetAllRoutes()
	api := app.Group("/api/v1")
	for i := 0; i < len(allRoutes); i++ {
		allRoutes[i](api)
	}
}
