package group_controllers

import (
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

func GetMyPosts(c *fiber.Ctx) error {
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.Preload("User").
		Where("user_id = ?", loggedInUserID).Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func GetPosts(c *fiber.Ctx) error {
	groupID := c.GetRespHeader("groupID")

	paginatedDB := API.Paginator(c)(initializers.DB)

	var posts []models.Post
	if err := paginatedDB.Preload("User").
		Where("group_id = ?", groupID).Find(&posts).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "",
		"posts":   posts,
	})
}

func AddPost(c *fiber.Ctx) error {
	parsedGroupID, err := uuid.Parse(c.GetRespHeader("groupID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Group ID."}
	}

	var reqBody schemas.PostCreateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Req Body"}
	}

	parsedLoggedInUserID, err := uuid.Parse(c.GetRespHeader("loggedInUserID"))
	if err != nil {
		return &fiber.Error{Code: 500, Message: "Error Parsing the Loggedin User ID."}
	}

	var user models.User
	if err := initializers.DB.Where("id=?", parsedLoggedInUserID).First(&user).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}
	if !user.IsVerified {
		return &fiber.Error{Code: 401, Message: config.VERIFICATION_ERROR}
	}

	if err := helpers.Validate[schemas.PostCreateSchema](reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	images, err := utils.UploadMultipleImages(c, "images", helpers.PostClient, 1280, 720)
	if err != nil {
		return err
	}

	newPost := models.Post{
		UserID:  parsedLoggedInUserID,
		Content: reqBody.Content,
		Images:  images,
		Tags:    reqBody.Tags,
		GroupID: parsedGroupID,
	}

	result := initializers.DB.Create(&newPost)
	if result.Error != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: result.Error.Error(), Err: result.Error}
	}

	routines.GetBlurHashesForPost(c, "images", &newPost)

	if err := initializers.DB.Preload("User").
		First(&newPost).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "Post Added",
		"post":    newPost,
	})
}

func UpdatePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").Preload("TaggedUsers").First(&post, "id = ? and user_id=?", parsedPostID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	var reqBody schemas.PostUpdateSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid Request Body."}
	}

	if reqBody.Content != "" {
		post.Content = reqBody.Content
	}
	if reqBody.Tags != nil {
		post.Tags = *reqBody.Tags
	}

	post.Edited = true

	if err := initializers.DB.Save(&post).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	postID := c.Params("postID")
	loggedInUserID := c.GetRespHeader("loggedInUserID")

	parsedPostID, err := uuid.Parse(postID)
	if err != nil {
		return &fiber.Error{Code: 400, Message: "Invalid ID"}
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, "id = ? AND user_id=?", parsedPostID, loggedInUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &fiber.Error{Code: 400, Message: "No Post of this ID found."}
		}
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	if err := initializers.DB.Delete(&post).Error; err != nil {
		return helpers.AppError{Code: 500, Message: config.DATABASE_ERROR, LogMessage: err.Error(), Err: err}
	}

	for _, image := range post.Images {
		go routines.DeleteFromBucket(helpers.PostClient, image)
	}

	return c.Status(204).JSON(fiber.Map{
		"status":  "success",
		"message": "Post deleted successfully",
	})
}
