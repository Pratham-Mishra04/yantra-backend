package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AnnouncementRouter(app *fiber.App) {
	announcementRoutes := app.Group("/announcement", middlewares.Protect)
	announcementRoutes.Get("/", group_controllers.GetAnnouncements)
	announcementRoutes.Post("/", group_controllers.AddAnnouncement)

	announcementRoutes.Patch("/:pageID", group_controllers.UpdateAnnouncement)
	announcementRoutes.Delete("/:pageID", group_controllers.DeleteAnnouncement)
}
