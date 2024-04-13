package routers

import (
	"github.com/Pratham-Mishra04/yantra-backend/controllers/messaging_controllers"
	"github.com/Pratham-Mishra04/yantra-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MessagingRouter(app *fiber.App) {
	messagingRoutes := app.Group("/messaging", middlewares.Protect)

	messagingRoutes.Get("/me", messaging_controllers.GetUserNonPopulatedChats)

	messagingRoutes.Get("/personal", messaging_controllers.GetPersonalChats)
	messagingRoutes.Get("/personal/unfiltered", messaging_controllers.GetPersonalUnFilteredChats)
	messagingRoutes.Get("/personal/unread", messaging_controllers.GetUnreadChats)
	messagingRoutes.Get("/group", messaging_controllers.GetGroupChats)

	messagingRoutes.Get("/:chatID", messaging_controllers.GetChat)
	messagingRoutes.Get("/group/:chatID", messaging_controllers.GetGroupChat)

	messagingRoutes.Get("/accept/:chatID", messaging_controllers.AcceptChat)

	messagingRoutes.Post("/chat", messaging_controllers.AddChat)
	messagingRoutes.Patch("/chat/last_read/:chatID", messaging_controllers.UpdateLastRead)

	messagingRoutes.Post("/chat/block", messaging_controllers.BlockChat)
	messagingRoutes.Post("/chat/unblock", messaging_controllers.UnblockChat)
	messagingRoutes.Post("/chat/reset", messaging_controllers.ResetChat)

	messagingRoutes.Get("/content/:chatID", messaging_controllers.GetMessages)
	messagingRoutes.Get("/content/group/:chatID", messaging_controllers.GetGroupChatMessages)

	messagingRoutes.Post("/content", messaging_controllers.AddMessage)
	messagingRoutes.Post("/content/group", messaging_controllers.AddGroupChatMessage)

	messagingRoutes.Delete("/content/:messageID", messaging_controllers.DeleteMessage)
	messagingRoutes.Delete("/content/project/:messageID", messaging_controllers.DeleteMessage)
}
