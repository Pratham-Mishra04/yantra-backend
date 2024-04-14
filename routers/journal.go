package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/user_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func JournalRouter(app *fiber.App) {
	journalRoutes := app.Group("/journal", middlewares.Protect)
	journalRoutes.Get("/", user_controllers.GetPages)
	journalRoutes.Post("/", user_controllers.CreatePage)

	journalRoutes.Patch("/:pageID", user_controllers.UpdatePage)
	journalRoutes.Delete("/:pageID", user_controllers.DeletePage)

	journalRoutes.Post("/ner", func(c *fiber.Ctx) error {
		var reqBody struct {
			Content string `json:"content"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		NERs := helpers.NERExtractionFromOnboarding(reqBody.Content)

		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"NERs":   NERs,
		})
	})

	journalRoutes.Post("/emotion", func(c *fiber.Ctx) error {
		var reqBody struct {
			Content string `json:"content"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		emotions, scores := helpers.EmotionExtractionFromOnboarding(reqBody.Content)

		return c.Status(200).JSON(fiber.Map{
			"status":   "success",
			"emotions": emotions,
			"scores":   scores,
		})
	})
}
