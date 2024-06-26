package main

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/Pratham-Mishra04/yantra-backend/routers"
	"github.com/Pratham-Mishra04/yantra-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.AddLogger()
	initializers.ConnectToCache()
	initializers.AutoMigrate()
	helpers.InitializeBucketClients()
	config.InitializeOAuthGoogle()

	utils.Repeater(func() {
		var groups []models.Group
		initializers.DB.
			Preload("Memberships").
			Preload("Memberships.User").
			Preload("Memberships.User.Journal").
			Find(&groups)

		for _, group := range groups {
			helpers.GroupDominatingEmotion(&group)
		}
	}, 500)

}

func main() {
	defer initializers.LoggerCleanUp()

	app := fiber.New(fiber.Config{
		ErrorHandler: helpers.ErrorHandler,
		BodyLimit:    config.BODY_LIMIT,
	})

	app.Use(helmet.New())
	app.Use(config.CORS())
	// app.Use(config.RATE_LIMITER())

	if initializers.CONFIG.ENV == initializers.DevelopmentEnv {
		app.Use(logger.New())
	}

	app.Use(logger.New())

	app.Static("/", "./public")

	routers.Config(app)

	app.Listen(":" + initializers.CONFIG.PORT)
}
