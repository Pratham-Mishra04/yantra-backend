package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ResourceRouter(app *fiber.App) {
	app.Get("resource/:groupID/serve/:resourceFileID", group_controllers.ServeResourceFile)

	resourceRoutes := app.Group("/resource", middlewares.Protect)

	resourceRoutes.Get("/", group_controllers.GetResourceBuckets)
	// resourceRoutes.Get("/:resourceBucketID", middlewares.BucketAuthorization("view"), group_controllers.GetResourceBucketFiles)
	resourceRoutes.Post("/", group_controllers.AddResourceBucket)
	resourceRoutes.Patch("/:resourceBucketID", group_controllers.EditResourceBucket)
	resourceRoutes.Delete("/:resourceBucketID", group_controllers.DeleteResourceBucket)

	resourceFileRoutes := resourceRoutes.Group("/:resourceBucketID/file")

	// resourceFileRoutes.Post("/", middlewares.BucketAuthorization("edit"), group_controllers.AddResourceFile)
	resourceFileRoutes.Patch("/:resourceFileID", group_controllers.EditResourceFile)
	resourceFileRoutes.Delete("/:resourceFileID", group_controllers.DeleteResourceFile)
}
