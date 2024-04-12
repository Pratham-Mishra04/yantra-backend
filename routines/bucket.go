package routines

import (
	"github.com/Pratham-Mishra04/yantra-backend/helpers"
	"github.com/Pratham-Mishra04/yantra-backend/initializers"
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
