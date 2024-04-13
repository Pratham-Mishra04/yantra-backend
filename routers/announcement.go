package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AnnouncementRouter(app *fiber.App) {
	announcementRoutes := app.Group("/announcement", middlewares.Protect, middlewares.AttachGroupHeader)
	announcementRoutes.Get("/", group_controllers.GetAnnouncements)
	announcementRoutes.Post("/", middlewares.ModeratorOnly, group_controllers.AddAnnouncement)

	announcementRoutes.Patch("/:announcementID", middlewares.ModeratorOnly, group_controllers.UpdateAnnouncement)
	announcementRoutes.Delete("/:announcementID", middlewares.ModeratorOnly, group_controllers.DeleteAnnouncement)
}
