package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CallbackRouter(app *fiber.App) {
	dyteCallbackRoutes := app.Group("/callbacks/dyte", middlewares.VerifyDyteWebHook)

	dyteCallbackRoutes.Post("/event", group_controllers.PostEventCallback)
}
