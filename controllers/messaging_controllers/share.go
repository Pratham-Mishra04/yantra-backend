package messaging_controllers

import (
	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func ShareItem(shareType string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		loggedInUserID := c.GetRespHeader("loggedInUserID")
		parsedUserID, _ := uuid.Parse(loggedInUserID)

		var reqBody struct {
			Content        string         `json:"content"`
			Chats          pq.StringArray `json:"chats"`
			PostID         string         `json:"postID"`
			EventID        string         `json:"eventID"`
			ProfileID      string         `json:"profileID"`
			AnnouncementID string         `json:"announcementID"`
		}
		if err := c.BodyParser(&reqBody); err != nil {
			return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
		}

		chats := reqBody.Chats

		for _, chatID := range chats {
			message := models.Message{
				UserID:  parsedUserID,
				Content: reqBody.Content,
			}

			parsedChatID, err := uuid.Parse(chatID)
			if err != nil {
				return &fiber.Error{Code: 400, Message: "Invalid ID."}
			}

			var chat models.Chat
			if err := initializers.DB.Where("id=?", parsedChatID).First(&chat).Error; err != nil {
				continue
			}

			if parsedUserID == chat.AcceptingUserID && chat.BlockedByCreatingUser {
				continue
			}

			if parsedUserID == chat.CreatingUserID && chat.BlockedByAcceptingUser {
				continue
			}

			message.ChatID = chat.ID

			switch shareType {
			case "post":
				parsedPostID, err := uuid.Parse(reqBody.PostID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Project ID."}
				}
				message.PostID = &parsedPostID
			case "announcement":
				parsedAnnouncementID, err := uuid.Parse(reqBody.AnnouncementID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Announcement ID."}
				}
				message.AnnouncementID = &parsedAnnouncementID
			case "event":
				parsedEventID, err := uuid.Parse(reqBody.EventID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Event ID."}
				}
				message.EventID = &parsedEventID
			case "profile":
				parsedProfileID, err := uuid.Parse(reqBody.ProfileID)
				if err != nil {
					return &fiber.Error{Code: 400, Message: "Invalid Profile ID."}
				}
				message.ProfileID = &parsedProfileID
			default:
				return &fiber.Error{Code: 400, Message: "Invalid Share Type."}
			}

			result := initializers.DB.Create(&message)
			if result.Error != nil {
				return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
			}
		}
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Shared",
		})

	}
}
