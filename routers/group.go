package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func GroupRouter(app *fiber.App) {
	groupRoutes := app.Group("/group", middlewares.Protect)
	groupRoutes.Get("/recommended", group_controllers.GetRecommendedGroups)

	groupRoutes.Get("/:groupID", group_controllers.GetGroup)
	//TODO only isModerator can access
	groupRoutes.Post("/", group_controllers.CreateGroup)

	groupRoutes.Post("/initial", group_controllers.JoinInitialGroup)
	groupRoutes.Post("/join/:groupID", group_controllers.JoinGroup)

	groupRoutes.Get("/feed", middlewares.AttachGroupHeader, group_controllers.GetCombinedFeed)

	groupRoutes.Patch("/", group_controllers.EditGroup)
	groupRoutes.Delete("/", group_controllers.DeleteGroup)
}
