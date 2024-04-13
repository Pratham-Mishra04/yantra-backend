package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ResourceRouter(app *fiber.App) {
	app.Get("resource/serve/:resourceFileID", middlewares.ResourceFileProtect, middlewares.AttachGroupHeader, group_controllers.ServeResourceFile)

	resourceRoutes := app.Group("/resource", middlewares.Protect, middlewares.AttachGroupHeader)

	resourceRoutes.Get("/", group_controllers.GetResourceBuckets)
	// resourceRoutes.Get("/:resourceBucketID", middlewares.BucketAuthorization("view"), group_controllers.GetResourceBucketFiles)
	resourceRoutes.Post("/", middlewares.ModeratorOnly, group_controllers.AddResourceBucket)
	resourceRoutes.Patch("/:resourceBucketID", middlewares.ModeratorOnly, group_controllers.EditResourceBucket)
	resourceRoutes.Delete("/:resourceBucketID", middlewares.ModeratorOnly, group_controllers.DeleteResourceBucket)

	resourceFileRoutes := resourceRoutes.Group("/:resourceBucketID/file")

	// resourceFileRoutes.Post("/", middlewares.BucketAuthorization("edit"), group_controllers.AddResourceFile)
	resourceFileRoutes.Patch("/:resourceFileID", group_controllers.EditResourceFile)
	resourceFileRoutes.Delete("/:resourceFileID", group_controllers.DeleteResourceFile)
}
