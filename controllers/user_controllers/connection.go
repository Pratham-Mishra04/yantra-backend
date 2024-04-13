package user_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPendingConnectionRequests(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var connections []models.Connection
	if err := initializers.DB.Preload("Sender").Preload("Receiver").
		Where("receiver_id = ? AND status = 0", userID).Find(&connections).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"connections": connections,
	})
}

func HandleConnectionRequest(handleType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

		connectionID := c.Get("connectionID")

		var connection models.Connection
		if err := initializers.DB.Where("id = ? AND receiver_id = ? AND status = 0", connectionID, userID).First(&connection).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return &fiber.Error{Code: 400, Message: "No Connection Request found."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		if handleType == "withdraw" {
			if err := initializers.DB.Delete(&connection).Error; err != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "connection request withdrawn",
			})
		}

		if handleType == "accept" {
			connection.Status = 1
		} else if handleType == "reject" {
			connection.Status = -1
		}

		if err := initializers.DB.Save(&connection).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Connection Request Handled",
		})
	}
}

func SendConnectionRequest(c *fiber.Ctx) error {
	senderID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var reqBody struct {
		UserID string `json:"userID"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	var user models.User
	if err := initializers.DB.Where("id = ?", reqBody.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No User of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	connection := models.Connection{
		SenderID:   senderID,
		ReceiverID: user.ID,
	}

	if err := initializers.DB.Create(&connection).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Connection Request Sent",
	})
}
