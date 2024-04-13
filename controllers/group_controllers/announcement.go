package group_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/Pratham-Mishra04/yantra-backend/schemas"
	API "github.com/Pratham-Mishra04/yantra-backend/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetAnnouncements(c *fiber.Ctx) error {
	parsedGroupID, err := uuid.Parse(c.Params("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var announcements []models.Announcement
	if err := paginatedDB.Where("group_id = ?", parsedGroupID).Order("created_at DESC").Find(&announcements).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"announcements": announcements,
	})
}

func AddAnnouncement(c *fiber.Ctx) error {
	var reqBody schemas.AnnouncementCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedGroupID, err := uuid.Parse(c.Params("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	announcement := models.Announcement{
		GroupID: parsedGroupID,
		Title:   reqBody.Title,
		Content: reqBody.Content,
	}

	if err := initializers.DB.Create(&announcement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":       "success",
		"message":      "Announcement Added",
		"announcement": announcement,
	})
}

func EditAnnouncement(c *fiber.Ctx) error {
	var reqBody schemas.AnnouncementUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedAnnouncementID, err := uuid.Parse(c.Params("announcementID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Announcement ID."}
	}

	parsedGroupID, err := uuid.Parse(c.Params("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	var announcement models.Announcement
	if err := initializers.DB.Where("id=? AND group_id = ?", parsedAnnouncementID, parsedGroupID).First(&announcement).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Announcement does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if reqBody.Title != "" {
		announcement.Title = reqBody.Title
	}
	if reqBody.Content != "" {
		announcement.Content = reqBody.Content
	}
	announcement.IsEdited = true

	if err := initializers.DB.Save(&announcement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":       "success",
		"message":      "Announcement Edited",
		"announcement": announcement,
	})
}

func DeleteAnnouncement(c *fiber.Ctx) error {
	parsedAnnouncementID, err := uuid.Parse(c.Params("announcementID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Announcement ID."}
	}

	parsedGroupID, err := uuid.Parse(c.Params("groupID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Group ID."}
	}

	var announcement models.Announcement
	if err := initializers.DB.Where("id=? AND group_id = ?", parsedAnnouncementID, parsedGroupID).First(&announcement).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Announcement does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&announcement).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Announcement Deleted",
	})
}
