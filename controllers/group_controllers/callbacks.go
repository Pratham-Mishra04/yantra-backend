package group_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MeetingEndedEvent struct {
	Event   string `json:"event"`
	Meeting struct {
		ID          string `json:"id"`
		SessionID   string `json:"sessionId"`
		Title       string `json:"title"`
		Status      string `json:"status"`
		CreatedAt   string `json:"createdAt"`
		StartedAt   string `json:"startedAt"`
		EndedAt     string `json:"endedAt"`
		OrganizedBy struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"organizedBy"`
	} `json:"meeting"`
	Reason string `json:"reason"`
}

func PostEventCallback(c *fiber.Ctx) error {
	var reqBody MeetingEndedEvent

	if err := c.BodyParser(&reqBody); err != nil {
		return &helpers.AppError{Code: 500, Message: "Failed to Handle Webhook", LogMessage: err.Error(), Err: err}
	}

	eventID := reqBody.Meeting.Title

	meetingID := reqBody.Meeting.ID

	users, err := helpers.GetDyteMeetingParticipants(meetingID)
	if err != nil {
		return err
	}

	var event models.Event
	if err := initializers.DB.First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, user := range users {
		review := models.Review{
			UserID:    user.ID,
			EventID:   &event.ID,
			Content:   "",
			IsPending: true,
		}

		result := initializers.DB.Save(&review)
		if result.Error != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Callback Handled.",
	})
}
