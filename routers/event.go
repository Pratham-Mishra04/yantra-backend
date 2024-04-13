package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func EventRouter(app *fiber.App) {
	eventRoutes := app.Group("/event", middlewares.Protect)
	eventRoutes.Get("/", group_controllers.GetEvents)
	eventRoutes.Post("/", group_controllers.AddAnnouncement)

	eventRoutes.Get("/:eventID", group_controllers.GetEvent)

	eventRoutes.Patch("/:pageID", group_controllers.UpdateAnnouncement)
	eventRoutes.Delete("/:pageID", group_controllers.DeleteAnnouncement)
}
