package group_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/yantra-backend/cache"
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/Pratham-Mishra04/yantra-backend/routines"
	"github.com/Pratham-Mishra04/yantra-backend/schemas"
	"github.com/Pratham-Mishra04/yantra-backend/utils"
	API "github.com/Pratham-Mishra04/yantra-backend/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	eventInCache, err := cache.GetEvent(eventID)
	if err == nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "",
			"event":   eventInCache,
		})
	}

	var event models.Event
	if err := initializers.DB.Preload("Group").Preload("Group.Moderator").Preload("Group.Moderator.User").Where("id = ?", eventID).First(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go cache.SetEvent(event.ID.String(), &event)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"event":   event,
	})
}

func GetEvents(c *fiber.Ctx) error {
	parsedGroupID, err := uuid.Parse(c.GetRespHeader("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	var events []models.Event
	if err := paginatedDB.Preload("Group").Preload("Group.Moderator").Preload("Group.Moderator.User").Where("group_id = ?", parsedGroupID).Order("created_at DESC").Find(&events).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"events": events,
	})
}

func AddEvent(c *fiber.Ctx) error {
	var reqBody schemas.EventCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	if err := helpers.Validate[schemas.EventCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	parsedGroupID, err := uuid.Parse(c.GetRespHeader("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	picName, err := utils.UploadImage(c, "coverPic", helpers.EventClient, 1920, 1080)
	if err != nil {
		return err
	}

	startTime, err := time.Parse(time.RFC3339, reqBody.StartTime)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Start Time."}
	}

	endTime, err := time.Parse(time.RFC3339, reqBody.EndTime)
	if err != nil || endTime.Before(startTime) {
		return &fiber.Error{Code: 400, Message: "Invalid End Time."}
	}

	event := models.Event{
		GroupID:     parsedGroupID,
		Title:       reqBody.Title,
		Tagline:     reqBody.Tagline,
		Description: reqBody.Description,
		Tags:        reqBody.Tags,
		Category:    reqBody.Category,
		Links:       reqBody.Links,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    reqBody.Location,
	}

	if picName != "" {
		event.CoverPic = picName
	}

	result := initializers.DB.Create(&event)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	err = helpers.CreateDyteMeeting(&event)
	if err != nil {
		return err
	}

	routines.GetImageBlurHash(c, "coverPic", &event)

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Event Added",
		"event":   event,
	})
}

func UpdateEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventID")

	var event models.Event
	if err := initializers.DB.Where("id = ?", eventID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody schemas.EventUpdateSchema
	c.BodyParser(&reqBody)

	picName, err := utils.UploadImage(c, "coverPic", helpers.EventClient, 1920, 1080)
	if err != nil {
		return err
	}
	oldEventPic := event.CoverPic

	if reqBody.Tagline != "" {
		event.Tagline = reqBody.Tagline
	}
	if picName != "" {
		event.CoverPic = picName
	}
	if reqBody.Category != "" {
		event.Category = reqBody.Category
	}
	if reqBody.Description != "" {
		event.Description = reqBody.Description
	}
	if reqBody.Location != "" {
		event.Location = reqBody.Location
	}
	if reqBody.Tags != nil {
		event.Tags = reqBody.Tags
	}
	if reqBody.Links != nil {
		event.Links = reqBody.Links
	}
	if reqBody.StartTime != "" {
		//TODO update on dyte
		startTime, err := time.Parse(time.RFC3339, reqBody.StartTime)
		if err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Start Time."}
		}
		event.StartTime = startTime
	}
	if reqBody.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, reqBody.EndTime)
		if err != nil || endTime.Before(event.StartTime) {
			return &fiber.Error{Code: 400, Message: "Invalid End Time."}
		}

		event.EndTime = endTime
	}

	if err := initializers.DB.Save(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	if reqBody.CoverPic != "" {
		go routines.DeleteFromBucket(helpers.EventClient, oldEventPic)
	}

	routines.GetImageBlurHash(c, "coverPic", &event)
	go cache.RemoveEvent(event.ID.String())

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Event updated successfully",
		"event":   event,
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	//TODO add OTP
	eventID := c.Params("eventID")

	parsedGroupID, err := uuid.Parse(c.GetRespHeader("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	var event models.Event
	if err := initializers.DB.Where("id = ? AND group_id=?", eventID, parsedGroupID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	eventPic := event.CoverPic

	if err := initializers.DB.Delete(&event).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	go routines.DeleteFromBucket(helpers.EventClient, eventPic)
	go cache.RemoveEvent(event.ID.String())

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}

func JoinLiveEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventID")
	parsedGroupID, err := uuid.Parse(c.GetRespHeader("groupID"))
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Group ID."}
	}

	var event models.Event
	if err := initializers.DB.Preload("Group").Preload("Group.Moderator").Where("id = ? AND group_id=?", eventID, parsedGroupID).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Event of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	var user models.User
	if err := initializers.DB.Where("id = ?", parsedUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	authToken, err := helpers.GetDyteMeetingAuthToken(&event, &user)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status":    "success",
		"message":   "",
		"authToken": authToken,
	})
}
