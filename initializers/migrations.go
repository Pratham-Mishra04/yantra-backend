package initializers

import (
	"fmt"

	"github.com/Pratham-Mishra04/yantra-backend/models"
)

func AutoMigrate() {
	fmt.Println("\nStarting Migrations...")
	DB.AutoMigrate(
		&models.User{},

		&models.UserVerification{},
		&models.OAuth{},
	)
	fmt.Println("Migrations Finished!")
}
