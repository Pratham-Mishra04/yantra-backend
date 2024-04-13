package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/user_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func JournalRouter(app *fiber.App) {
	journalRoutes := app.Group("/journal", middlewares.Protect)
	journalRoutes.Get("/", user_controllers.GetPages)
	journalRoutes.Post("/", user_controllers.CreatePage)

	journalRoutes.Patch("/:pageID", user_controllers.UpdatePage)
	journalRoutes.Delete("/:pageID", user_controllers.DeletePage)
}
