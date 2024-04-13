package user_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	API "github.com/Pratham-Mishra04/yantra-backend/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPages(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var pages []models.Page
	if err := paginatedDB.Preload("User").
		Where("user_id = ?", userID).Find(&pages).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"pages":   pages,
	})
}

func CreatePage(c *fiber.Ctx) error {
	var reqBody struct {
		Content string `json:"content" validate:"required,max=2500"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	userID := c.GetRespHeader("loggedInUserID")

	var journal models.Journal
	if err := initializers.DB.First(&journal, "user_id=?", userID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	page := models.Page{
		Content:   reqBody.Content,
		JournalID: journal.ID,
	}

	if err := initializers.DB.Create(&page).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Page created!",
	})
}

func UpdatePage(c *fiber.Ctx) error {
	pageID := c.Params("pageID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var page models.Page
	if err := initializers.DB.First(&page, "id = ? and user_id=?", pageID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Page of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody struct {
		Content string `json:"content" validate:"required,max=2500"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if err := initializers.DB.Save(&page).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Page updated successfully",
	})
}

func DeletePage(c *fiber.Ctx) error {
	pageID := c.Params("pageID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	var page models.Page
	if err := initializers.DB.Preload("User").First(&page, "id = ? AND user_id=?", pageID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Page of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&page).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Page deleted successfully",
	})
}
