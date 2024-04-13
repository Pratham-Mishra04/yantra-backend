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

func BucketAuthorization(action string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		groupID := c.GetRespHeader("groupID")
		resourceBucketID := c.Params("resourceBucketID")

		var group models.Group
		if err := initializers.DB.Preload("Moderator").Where("id = ?", groupID).First(&group).Error; err != nil {
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		var resourceBucket models.ResourceBucket

		if err := initializers.DB.Where("id=? AND group_id = ?", resourceBucketID, groupID).First(&resourceBucket).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Resource Bucket does not exist."}
			}
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}

		if loggedInUserID == group.Moderator.UserID.String() {
			return c.Next()
		}

		check := false

		if action == "view" && !resourceBucket.OnlyAdminViewAccess {
			check = true
		} else if action == "edit" && !resourceBucket.OnlyAdminEditAccess {
			check = true
		}

		if !check {
			return &fiber.Error{Code: 403, Message: "Cannot access this Resource Bucket."}
		}

		return c.Next()
	}
}
