package routers

import (
	"github.com/gofiber/fiber/v2"
)

func Config(app *fiber.App) {
	AnnouncementRouter(app)
	CallbackRouter(app)
	ConnectionRouter(app)
	EventRouter(app)
	GroupRouter(app)
	JournalRouter(app)
	MessagingRouter(app)
	OauthRouter(app)
	PollRouter(app)
	PostRouter(app)
	ResourceRouter(app)
	ReviewRouter(app)
	UserRouter(app)
	VerificationRouter(app)
}
