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

func GetPendingReviews(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var reviews []models.Review
	if err := initializers.DB.Where("user_id = ? AND is_pending = ?", userID, true).Find(&reviews).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"reviews": reviews,
	})
}

func AddReview(c *fiber.Ctx) error {
	var reqBody struct {
		Content string `json:"content" validate:"required,max=16"`
		EventID string `json:"eventID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedLoggedInUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	review := models.Review{
		UserID:    parsedLoggedInUserID,
		Content:   reqBody.Content,
		IsPending: false,
	}

	if reqBody.EventID != "" {
		var event models.Event
		if err := initializers.DB.Where("id=?", reqBody.EventID).First(&event).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Event does not exist."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		review.EventID = &event.ID
	}

	result := initializers.DB.Create(&review)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Review Added",
	})
}
