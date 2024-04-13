package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/user_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ReviewRouter(app *fiber.App) {
	reviewRoutes := app.Group("/review", middlewares.Protect)
	reviewRoutes.Get("/", user_controllers.GetPendingReviews)
	reviewRoutes.Post("/", user_controllers.AddReview)
}
