package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func GroupRouter(app *fiber.App) {
	groupRoutes := app.Group("/group", middlewares.Protect)
	groupRoutes.Get("/:groupID", group_controllers.GetGroup)
	groupRoutes.Post("/", group_controllers.CreateGroup)

	groupRoutes.Patch("/", group_controllers.EditGroup)
	groupRoutes.Delete("/", group_controllers.DeleteGroup)
}
