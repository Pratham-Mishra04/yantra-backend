package initializers

import (
	"fmt"

	"github.com/Pratham-Mishra04/yantra-backend/models"
)

func AutoMigrate() {
	fmt.Println("\nStarting Migrations...")
	DB.AutoMigrate(
		&models.User{},
		&models.Announcement{},
		&models.Comment{},
		&models.Connection{},
		&models.Event{},
		&models.Group{},
		&models.Journal{},
		&models.Notification{},
		&models.Option{},
		&models.Page{},
		&models.Poll{},
		&models.Post{},
		&models.Report{},
		&models.ResourceBucket{},
		&models.ResourceFile{},
		&models.Review{},

		&models.UserVerification{},
		&models.OAuth{},
	)
	fmt.Println("Migrations Finished!")
}
