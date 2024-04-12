package routines

import (
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
	"github.com/Pratham-Mishra04/yantra-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DeleteFromBucket(client *helpers.BucketClient, path string) {
	if path == "" || path == "default.jpg" {
		return
	}
	err := client.DeleteBucketFile(path)
	if err != nil {
		initializers.Logger.Warnw("Error while deleting file from bucket", "Error", err)
	}
}

func IncrementResourceBucketFiles(resourceBucketID uuid.UUID) {
	var resourceBucket models.ResourceBucket
	if err := initializers.DB.First(&resourceBucket, "id = ?", resourceBucketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Resource Bucket of this ID found-IncrementResourceBucketFiles.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	} else {
		resourceBucket.NumberOfFiles++
		if err := initializers.DB.Save(&resourceBucket).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	}
}

func DecrementResourceBucketFiles(resourceBucketID uuid.UUID) {
	var resourceBucket models.ResourceBucket
	if err := initializers.DB.First(&resourceBucket, "id = ?", resourceBucketID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.LogDatabaseError("No Resource Bucket of this ID found-IncrementResourceBucketFiles.", err, "go_routine")
		} else {
			helpers.LogDatabaseError("Error while fetching Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	} else {
		resourceBucket.NumberOfFiles--
		if err := initializers.DB.Save(&resourceBucket).Error; err != nil {
			helpers.LogDatabaseError("Error while updating Resource Bucket-IncrementResourceBucketFiles", err, "go_routine")
		}
	}
}
