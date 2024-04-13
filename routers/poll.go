package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/group_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func PollRouter(app *fiber.App) {
	pollRouter := app.Group("/poll", middlewares.Protect)
	pollRouter.Get("/", group_controllers.GetPolls)

	pollRouter.Post("/", group_controllers.CreatePoll)
	pollRouter.Patch("/:pollID", group_controllers.EditPoll)
	pollRouter.Delete("/:pollID", group_controllers.DeletePoll)

	pollRouter.Patch("/vote/:pollID/:OptionID", group_controllers.VotePoll)
	pollRouter.Patch("/unvote/:pollID/:OptionID", group_controllers.UnvotePoll)
}
