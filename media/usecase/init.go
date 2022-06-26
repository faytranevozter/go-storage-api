package usecase

import (
	"storage-api/domain"
	"time"
)

type uMedia struct {
	repoMongo      domain.MongoRepository
	repoGBucket    domain.GoogleBucketRepository
	repoCloudinary domain.CloudinaryRepository
	contextTimeout time.Duration
}

// NewUsecaseMedia ...
func NewUsecaseMedia(rm domain.MongoRepository, gbr domain.GoogleBucketRepository, cr domain.CloudinaryRepository, timeout time.Duration) domain.MediaUsecase {
	return &uMedia{
		repoMongo:      rm,
		repoGBucket:    gbr,
		repoCloudinary: cr,
		contextTimeout: timeout,
	}
}
