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
		&models.GroupMembership{},
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

		&models.Chat{},
		&models.Message{},
		&models.GroupChat{},
		&models.GroupChatMessage{},

		&models.UserVerification{},
		&models.OAuth{},
	)
	fmt.Println("Migrations Finished!")
}
