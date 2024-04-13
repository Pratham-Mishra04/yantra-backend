package group_controllers

import (
	"sort"
	"time"

	"github.com/Pratham-Mishra04/yantra-backend/config"
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/Pratham-Mishra04/yantra-backend/utils"
	API "github.com/Pratham-Mishra04/yantra-backend/utils/APIFeatures"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CombinedFeedItem interface {
	GetCreatedAt() time.Time
}

type PostAlias models.Post

func (p PostAlias) GetCreatedAt() time.Time {
	return p.CreatedAt
}

type AnnouncementAlias models.Announcement

func (a AnnouncementAlias) GetCreatedAt() time.Time {
	return a.CreatedAt
}

type PollAlias models.Poll

func (p PollAlias) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func GetCombinedFeed(c *fiber.Ctx) error {
	parsedGroupID, _ := uuid.Parse(c.GetRespHeader("groupID"))

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select(utils.UserSelect)
		}).
		Preload("Group").
		Preload("Group.Moderator").
		Preload("Group.Moderator.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(utils.UserSelect)
		}).
		Joins("JOIN users ON posts.user_id = users.id AND users.active = ?", true).
		Where("group_id = ?", parsedGroupID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	paginatedDB = API.Paginator(c)(initializers.DB)

	var announcements []models.Announcement
	if err := paginatedDB.
		Preload("Group").
		Preload("Group.Moderator").
		Preload("Group.Moderator.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(utils.UserSelect)
		}).
		Where("group_id = ?", parsedGroupID).
		Order("created_at DESC").
		Find(&announcements).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	paginatedDB = API.Paginator(c)(initializers.DB)

	db := paginatedDB.
		Preload("Group").
		Preload("Group.Moderator").
		Preload("Group.Moderator.User", func(db *gorm.DB) *gorm.DB {
			return db.Select(utils.UserSelect)
		}).
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.created_at DESC")
		}).
		Preload("Options.VotedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select(utils.UserSelect).Limit(3)
		}).
		Where("group_id = ?", parsedGroupID)

	var polls []models.Poll
	if err := db.Order("created_at DESC").Find(&polls).Error; err != nil {
		return helpers.AppError{Code: fiber.StatusInternalServerError, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var combinedFeed []CombinedFeedItem
	for _, a := range announcements {
		combinedFeed = append(combinedFeed, AnnouncementAlias(a))
	}
	for _, p := range polls {
		combinedFeed = append(combinedFeed, PollAlias(p))
	}
	for _, p := range posts {
		combinedFeed = append(combinedFeed, PostAlias(p))
	}

	sort.Slice(combinedFeed, func(i, j int) bool {
		return combinedFeed[i].GetCreatedAt().After(combinedFeed[j].GetCreatedAt())
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"feed":   combinedFeed,
	})
}
