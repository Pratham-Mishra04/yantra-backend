package cache

import "github.com/Pratham-Mishra04/yantra-backend/models"

func GetResourceBucket(slug string) (*models.ResourceBucket, error) {
	var resourceBucket models.ResourceBucket
	err := GetFromCacheGeneric("resource_bucket-"+slug, &resourceBucket)
	return &resourceBucket, err
}

func SetResourceBucket(slug string, resourceBucket *models.ResourceBucket) error {
	return SetToCacheGeneric("resource_bucket-"+slug, resourceBucket)
}

func RemoveResourceBucket(slug string) error {
	return RemoveFromCacheGeneric("resource_bucket-" + slug)
}
