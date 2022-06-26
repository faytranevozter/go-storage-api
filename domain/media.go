package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AvailableUploadType List of available type document to upload
var AvailableUploadType []string = []string{
	"document",
	"image",
	"video",
	"audio",
	"pdf",
}

// AllowedMimeTypes List of available type document to upload
var AllowedMimeTypes []string = []string{
	"application/pdf",                                                                               // pdf
	"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", // doc & docx
	"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp", // images
	"video/mp4", "video/mpeg", "video/ogg", "video/webm", // videos
	"audio/mpeg", "audio/ogg", "audio/wav", "audio/webm", // auidos
}

// AllowedSort Allowed field to sort in mongo
var AllowedSort = []string{
	"type",
	"title",
	"created_at",
	"updated_at",
}

// Media data structure
type Media struct {
	ID                  primitive.ObjectID `json:"id"`
	BrandownerID        int64              `json:"brandowner_id"`
	Title               string             `json:"title"`
	Description         string             `json:"description"`
	Type                string             `json:"type"`
	Provider            string             `json:"provider"`
	PublicURL           string             `json:"public_url"`
	GoogleBucketStorage GoogleBucket       `json:"-"`
	CloudinaryStorage   CloudinaryData     `json:"-"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	DeletedAt           *time.Time         `json:"-"`
}

// MediaMongo data structure
type MediaMongo struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id"`
	BrandownerID        int64              `json:"brandowner_id" bson:"brandowner_id"`
	Title               string             `json:"title" bson:"title"`
	Description         string             `json:"description" bson:"description"`
	Type                string             `json:"type" bson:"type"`
	Provider            string             `json:"provider" bson:"provider"`
	PublicURL           string             `json:"public_url" bson:"public_url"`
	GoogleBucketStorage GoogleBucket       `json:"google_bucket" bson:"google_bucket"`
	CloudinaryStorage   CloudinaryData     `json:"cloudinary_storage" bson:"cloudinary_storage"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt           *time.Time         `json:"-" bson:"deleted_at"`
}

// GoogleBucket data structure
type GoogleBucket struct {
	BucketName string `json:"bucket_name" bson:"bucket_name"`
	ID         string `json:"id" bson:"id"`
	Name       string `json:"name" bson:"name"`
	Generation int64  `json:"generation" bson:"generation"`
}

// GoogleBucket data structure
type CloudinaryData struct {
	AssetID      string `json:"asset_id" bson:"asset_id"`
	PublicID     string `json:"public_id" bson:"public_id"`
	Version      int    `json:"version" bson:"version"`
	VersionID    string `json:"version_id" bson:"version_id"`
	Signature    string `json:"signature" bson:"signature"`
	Width        int    `json:"width" bson:"width"`
	Height       int    `json:"height" bson:"height"`
	Format       string `json:"format" bson:"format"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
}
