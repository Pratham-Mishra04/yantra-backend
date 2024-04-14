package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func PostRouter(app *fiber.App) {
	postRoutes := app.Group("/post", middlewares.Protect, middlewares.AttachGroupHeader)
	postRoutes.Get("/", group_controllers.GetPosts)
	postRoutes.Post("/", group_controllers.AddPost)

	postRoutes.Get("/my", group_controllers.GetMyPosts)

	postRoutes.Patch("/:postID", middlewares.ModeratorOnly, group_controllers.UpdatePost)
	postRoutes.Delete("/:postID", middlewares.ModeratorOnly, group_controllers.DeletePost)
}
