package routes

import (
	"grf/core/server"
)

func RegisterRoutes(app *server.App) {
	apiV1 := app.FiberApp.Group("/v1")
	RegisterAuthRoutes(apiV1, app)
}
