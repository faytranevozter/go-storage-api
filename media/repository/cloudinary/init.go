package cloudinaryrepo

import (
	"context"
	"mime/multipart"
	"os"
	"storage-api/domain"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	cloudinaryUploader "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryRepo struct {
	enabled       bool
	cloudinaryURL string
}

// NewCloudinaryRepo ...
func NewCloudinaryRepo() domain.CloudinaryRepository {
	var cloudinaryURL = os.Getenv("CLOUDINARY_URL")
	enabled, _ := strconv.ParseBool(os.Getenv("CLOUDINARY_ENABLED"))
	return &cloudinaryRepo{
		cloudinaryURL: cloudinaryURL,
		enabled:       enabled,
	}
}

func (r *cloudinaryRepo) IsEnabled() bool {
	return r.enabled
}

func (r *cloudinaryRepo) SupportedType() []string {
	return []string{
		"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp", // images
		"video/mp4", "video/mpeg", "video/ogg", "video/webm", // videos
	}
}

// UploadFile uploads an object.
func (r *cloudinaryRepo) UploadFile(ctx context.Context, f multipart.File, dir, objectName string) (*cloudinaryUploader.UploadResult, error) {
	cld, _ := cloudinary.NewFromURL(r.cloudinaryURL)

	return cld.Upload.Upload(ctx, f, cloudinaryUploader.UploadParams{
		PublicID: objectName,
		Folder:   dir,
	})
}
