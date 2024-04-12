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

func FetchPolls(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("orgID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid organization ID."}
	}

	paginatedDB := API.Paginator(c)(initializers.DB)

	db := paginatedDB.Preload("Options", func(db *gorm.DB) *gorm.DB {
		return db.Order("options.created_at DESC")
	}).Preload("Options.VotedBy", LimitedUsers).Where("organization_id = ?", orgID)

	var polls []models.Poll
	if err := db.Order("created_at DESC").Find(&polls).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"polls":  polls,
	})
}

func CreatePoll(c *fiber.Ctx) error {
	var reqBody schemas.CreatePollRequest

	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}
	if len(reqBody.Options) < 2 || len(reqBody.Options) > 10 {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	groupID, _ := uuid.Parse(c.Params("groupID"))

	var poll = models.Poll{
		GroupID:       groupID,
		Title:         reqBody.Title,
		Content:       reqBody.Content,
		IsMultiAnswer: reqBody.IsMultiAnswer,
	}

	if err := initializers.DB.Create(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()

	for _, optionText := range reqBody.Options {
		option := &models.Option{
			PollID:  poll.ID,
			Content: optionText,
		}
		if err := tx.Create(&option).Error; err != nil {
			tx.Rollback()
			return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Preload("Options").First(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll created!",
		"poll":    poll,
	})
}

func VotePoll(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	parsedPollID, err := uuid.Parse(c.Params("pollID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Poll ID."}
	}

	parsedOptionID, err := uuid.Parse(c.Params("OptionID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Option ID."}
	}

	var poll models.Poll
	if err := initializers.DB.Preload("Options").Preload("Options.VotedBy").First(&poll, "id = ?", parsedPollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	votedOptionID := uuid.Nil

	for _, option := range poll.Options {
		for _, voter := range option.VotedBy {
			if voter.ID == parsedUserID {
				votedOptionID = option.ID
				if !poll.IsMultiAnswer {
					return &fiber.Error{Code: fiber.StatusBadRequest, Message: "You have already voted"}
				}
			}
		}
	}

	var option models.Option
	if err := initializers.DB.First(&option, "id = ?", parsedOptionID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if votedOptionID == option.ID {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Vote recorded!",
		})
	}

	var user models.User
	if err := initializers.DB.First(&user, "id = ?", parsedUserID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	option.Votes++
	poll.TotalVotes++
	option.VotedBy = append(option.VotedBy, user)

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "VotePoll")
		}
	}()

	if err := tx.Save(&option).Error; err != nil {
		return err
	}

	if err := tx.Save(&poll).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Vote recorded!",
	})
}

func UnvotePoll(c *fiber.Ctx) error {
	parsedUserID, _ := uuid.Parse(c.GetRespHeader("loggedInUserID"))

	parsedOptionID, err := uuid.Parse(c.Params("OptionID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid Option ID."}
	}

	var option models.Option
	if err := initializers.DB.Preload("VotedBy").First(&option, "id = ?", parsedOptionID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var poll models.Poll
	if err := initializers.DB.First(&poll, "id = ?", option.PollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	isFound := false

	for i, voter := range option.VotedBy {
		if voter.ID == parsedUserID {
			option.VotedBy = append(option.VotedBy[:i], option.VotedBy[i+1:]...)
			option.Votes--
			poll.TotalVotes--
			isFound = true
			break
		}
	}

	if !isFound {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User has not voted"}
	}

	tx := initializers.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
			go helpers.LogDatabaseError("Transaction rolled back due to error", tx.Error, "UnVotePoll")
		}
	}()

	if err := tx.Model(&option).Association("VotedBy").Replace(option.VotedBy); err != nil {
		return err
	}

	if err := tx.Save(&option).Error; err != nil {
		return err
	}

	if err := tx.Save(&poll).Error; err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Vote removed!",
	})
}

func EditPoll(c *fiber.Ctx) error {
	var reqBody schemas.EditPollRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid request body."}
	}

	pollID, _ := uuid.Parse(c.Params("pollID"))

	var poll models.Poll
	if err := initializers.DB.First(&poll, "id = ?", pollID).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	poll.Content = reqBody.Content
	poll.IsEdited = true

	if err := initializers.DB.Save(&poll).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll edited!",
	})
}

func LimitedUsers(db *gorm.DB) *gorm.DB {
	return db.Limit(3)
}

func DeletePoll(c *fiber.Ctx) error {
	pollID, err := uuid.Parse(c.Params("pollID"))
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid poll ID."}
	}

	var poll models.Poll
	if err := initializers.DB.Preload("Options").First(&poll, "id = ?", pollID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No Poll if this ID Found."}
		}
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, option := range poll.Options {
		if err := tx.Model(&option).Association("VotedBy").Clear(); err != nil {
			tx.Rollback()
			return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
		}
	}

	if err := tx.Delete(&poll).Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Poll deleted!",
	})
}
