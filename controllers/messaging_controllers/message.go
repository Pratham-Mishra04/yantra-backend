package messaging_controllers

import (
	"time"

	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/Pratham-Mishra04/yantra-backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedLoggedInUserID, _ := uuid.Parse(loggedInUserID)

	var chat models.Chat
	if err := initializers.DB.Preload("LastReadMessageByAcceptingUser").
		Preload("LastReadMessageByCreatingUser").Where("id = ? AND (creating_user_id = ? OR accepting_user_id = ?)", chatID, parsedLoggedInUserID, parsedLoggedInUserID).
		First(&chat).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	timestamp := time.Now()
	if parsedLoggedInUserID == chat.AcceptingUserID {
		timestamp = chat.LastResetByAcceptingUser
	} else if parsedLoggedInUserID == chat.CreatingUserID {
		timestamp = chat.LastResetByCreatingUser
	}

	// paginatedDB := API.Paginator(c)(initializers.DB)

	var messages []models.Message
	// if err := paginatedDB.
	if err := initializers.DB.
		Preload("User").
		Preload("Post").
		Preload("Profile").
		Preload("Announcement").
		Preload("Post.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(utils.UserSelect)
		}).
		Where("chat_id = ? AND created_at > ?", chatID, timestamp).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	// if len(messages) > 0 {
	// 	go routines.UpdateChatLastRead(chat.ID, messages, parsedLoggedInUserID)
	// }

	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"messages": messages,
	})
}

func GetGroupChatMessages(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	var messages []models.GroupChatMessage
	if err := initializers.DB.
		Preload("User").
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Find(&messages).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"message":  "",
		"messages": messages,
	})
}

func AddMessage(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content string `json:"content"`
		ChatID  string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	var chat models.Chat
	if err := initializers.DB.First(&chat, "id=? AND (creating_user_id = ? OR accepting_user_id = ?)", parsedChatID, parsedUserID, parsedUserID).Error; err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	if parsedUserID == chat.CreatingUserID && chat.BlockedByAcceptingUser {
		return &fiber.Error{Code: 400, Message: "You have been blocked."}
	}

	if parsedUserID == chat.AcceptingUserID && chat.BlockedByCreatingUser {
		return &fiber.Error{Code: 400, Message: "You have been blocked."}
	}

	message := models.Message{
		UserID:  parsedUserID,
		Content: reqBody.Content,
		ChatID:  parsedChatID,
	}

	result := initializers.DB.Create(&message)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	chat.LatestMessageID = &message.ID

	result = initializers.DB.Save(&chat)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}

func AddGroupChatMessage(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")
	parsedUserID, _ := uuid.Parse(loggedInUserID)

	var reqBody struct {
		Content string `json:"content"`
		ChatID  string `json:"chatID"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	chatID := reqBody.ChatID

	message := models.GroupChatMessage{
		UserID:  parsedUserID,
		Content: reqBody.Content,
	}

	parsedChatID, err := uuid.Parse(chatID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID."}
	}

	message.ChatID = parsedChatID

	result := initializers.DB.Create(&message)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}

func DeleteMessage(c *fiber.Ctx) error {
	messageID := c.Params("messageID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedMessageID, err := uuid.Parse(messageID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var message models.Message
	if err := initializers.DB.First(&message, "id = ? AND user_id=?", parsedMessageID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Message of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	result := initializers.DB.Delete(&message)

	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Message Deleted",
	})
}
