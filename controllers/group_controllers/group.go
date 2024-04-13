package group_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/Pratham-Mishra04/yantra-backend/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetGroup(c *fiber.Ctx) error {
	groupID := c.Params("groupID")

	// groupInCache, err := cache.GetGroup(groupID)
	// if err == nil {
	// 	return c.Status(200).JSON(fiber.Map{
	// 		"status":  "success",
	// 		"message": "",
	// 		"group":   groupInCache,
	// 	})
	// }

	var group models.Group
	if err := initializers.DB.Where("id = ?", groupID).First(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// go cache.SetGroup(group.ID.String(), &group)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"group":   group,
	})
}

func CreateGroup(c *fiber.Ctx) error {
	var reqBody schemas.AnnouncementCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	group := models.Group{
		Title: reqBody.Title,
	}

	if err := initializers.DB.Create(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Group Created",
		"group":   group,
	})
}

func EditGroup(c *fiber.Ctx) error {
	return nil
}

func DeleteGroup(c *fiber.Ctx) error {
	parsedGroupID, err := uuid.Parse(c.Params("groupID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Group ID."}
	}

	var group models.Group
	if err := initializers.DB.Where("id = ?", parsedGroupID).First(&group).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Group does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Group Deleted",
	})
}
