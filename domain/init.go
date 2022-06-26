package domain

import (
	"context"
	"mime/multipart"

	"cloud.google.com/go/storage"
	cloudinaryUploader "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	PROVIDER_CLOUDINARY    string = "cloudinary"
	PROVIDER_GOOGLE_BUCKET string = "google_bucket"
)

var AVAILABLE_PROVIDER = []string{PROVIDER_CLOUDINARY, PROVIDER_GOOGLE_BUCKET}

// MediaUsecase represent media's usecase
type MediaUsecase interface {
	ListMedia(ctx context.Context, options DefaultPayload) BaseResponse
	DetailMedia(ctx context.Context, options DefaultPayload) BaseResponse
	UploadMedia(ctx context.Context, options DefaultPayload) BaseResponse
	UpdateMedia(ctx context.Context, options DefaultPayload) BaseResponse
	DeleteMedia(ctx context.Context, options DefaultPayload) BaseResponse
}

// MongoRepository represent the media's repository contract
type MongoRepository interface {
	FetchMedia(ctx context.Context, options map[string]interface{}) ([]MediaMongo, int64, error)
	InsertMedia(ctx context.Context, mMongo *MediaMongo) error
	UpdateMedia(ctx context.Context, mMongo *MediaMongo) error
}

// GoogleBucketRepository represent the media's repository contract
type GoogleBucketRepository interface {
	IsEnabled() bool
	SupportedType() []string
	UploadFile(ctx context.Context, file multipart.File, objectName string) (*storage.ObjectAttrs, error)
}

// CloudinaryRepository represent the media's repository contract
type CloudinaryRepository interface {
	IsEnabled() bool
	SupportedType() []string
	UploadFile(ctx context.Context, file multipart.File, directory, objectName string) (*cloudinaryUploader.UploadResult, error)
}
