package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/user_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ConnectionRouter(app *fiber.App) {
	connectionRoutes := app.Group("/connection", middlewares.Protect)
	connectionRoutes.Get("/", user_controllers.GetPendingConnectionRequests)
	connectionRoutes.Post("/", user_controllers.SendConnectionRequest)

	connectionRoutes.Patch("/:connectionID/accept", user_controllers.HandleConnectionRequest("accept"))
	connectionRoutes.Patch("/:connectionID/reject", user_controllers.HandleConnectionRequest("reject"))

	connectionRoutes.Delete("/:connectionID", user_controllers.HandleConnectionRequest("withdraw"))
}
