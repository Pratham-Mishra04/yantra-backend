package messaging_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetGroupChat(c *fiber.Ctx) error {
	chatID := c.Params("chatID")

	var chat models.GroupChat
	err := initializers.DB.
		Preload("User").
		Where("id = ?", chatID).
		First(&chat).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No Chat of this ID found."}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"chat":    chat,
	})
}

func EditGroupChat(c *fiber.Ctx) error {
	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	groupChatID := c.Params("chatID")

	var groupChat models.GroupChat
	err := initializers.DB.First(&groupChat, "id = ?", groupChatID).Error
	if err != nil {
		return &fiber.Error{Code: 400, Message: "No chat of this id found."}
	}

	if reqBody.Title != "" {
		groupChat.Title = reqBody.Title
	}
	if reqBody.Description != "" {
		groupChat.Description = reqBody.Description
	}

	// picName, err := utils.UploadImage(c, "coverPic", helpers.ChatClient, 720, 720)
	// if err != nil {
	// 	return err
	// }

	// oldGroupPic := groupChat.CoverPic

	// if picName != "" {
	// 	groupChat.CoverPic = picName
	// }

	result := initializers.DB.Save(&groupChat)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	// if picName != "" {
	// 	go routines.DeleteFromBucket(helpers.ChatClient, oldGroupPic)
	// }

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Chat Updated",
		"chat":    groupChat,
	})
}
