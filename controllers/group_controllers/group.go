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

func JoinInitialGroup(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var group models.Group
	if err := initializers.DB.First(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	membership := models.GroupMembership{
		UserID:  parsedUserID,
		GroupID: group.ID,
	}

	if err := initializers.DB.Create(&membership).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "",
		"group":      group,
		"membership": membership,
	})
}

func GetRecommendedGroups(c *fiber.Ctx) error {
	userID := c.GetRespHeader("loggedInUserID")

	var user models.User
	if err := initializers.DB.Preload("Journal").Where("id=?", userID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	groups := helpers.GetGroupRecommendations(&user)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"groups":  groups,
	})
}

func GetGroup(c *fiber.Ctx) error {
	groupID := c.GetRespHeader("groupID")

	// groupInCache, err := cache.GetGroup(groupID)
	// if err == nil {
	// 	return c.Status(200).JSON(fiber.Map{
	// 		"status":  "success",
	// 		"message": "",
	// 		"group":   groupInCache,
	// 	})
	// }

	var group models.Group
	if err := initializers.DB.Preload("Moderator").Preload("Moderator.User").Preload("Memberships").Preload("Memberships.User").Where("id = ?", groupID).First(&group).Error; err != nil {
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
	var reqBody schemas.GroupCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var moderator models.Moderator
	if err := initializers.DB.Where("user_id = ?", parsedUserID).First(&moderator).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	group := models.Group{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		ModeratorID: moderator.ID,
	}

	if err := initializers.DB.Create(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	membership := models.GroupMembership{
		UserID:  parsedUserID,
		GroupID: group.ID,
	}

	if err := initializers.DB.Create(&membership).Error; err != nil {
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

func JoinGroup(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	parsedGroupID, err := uuid.Parse(c.Params("groupID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Group ID."}
	}

	var group models.Group
	if err := initializers.DB.Where("id = ?", parsedGroupID).First(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	membership := models.GroupMembership{
		UserID:  parsedUserID,
		GroupID: group.ID,
	}

	//TODO increase no of members of the group.

	if err := initializers.DB.Create(&membership).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	group.NumberOfMembers++

	if err := initializers.DB.Save(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":     "success",
		"message":    "",
		"group":      group,
		"membership": membership,
	})
}
