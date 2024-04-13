package middlewares

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AttachGroupHeader(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var membership models.GroupMembership
	if err := initializers.DB.Where("user_id = ?", parsedUserID).First(&membership).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Membership does not exist."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	c.Set("groupID", membership.GroupID.String())

	return c.Next()
}

func ModeratorOnly(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	parsedGroupID, _ := uuid.Parse(c.GetRespHeader("groupID"))

	var group models.Group
	if err := initializers.DB.Preload("Moderator").Where("id = ?", parsedGroupID).First(&group).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if group.Moderator.UserID != parsedUserID {
		return &fiber.Error{Code: 401, Message: "Cannot Perform this action."}
	}

	return c.Next()
}
