package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func EventRouter(app *fiber.App) {
	eventRoutes := app.Group("/event", middlewares.Protect, middlewares.AttachGroupHeader)
	eventRoutes.Get("/", group_controllers.GetEvents)
	eventRoutes.Post("/", middlewares.ModeratorOnly, group_controllers.AddEvent)

	eventRoutes.Get("/token/:eventID", group_controllers.JoinLiveEvent)
	eventRoutes.Get("/:eventID", group_controllers.GetEvent)

	eventRoutes.Patch("/:eventID", middlewares.ModeratorOnly, group_controllers.UpdateEvent)
	eventRoutes.Delete("/:eventID", middlewares.ModeratorOnly, group_controllers.DeleteEvent)
}
