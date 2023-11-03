package googlebucket

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"storage-api/domain"
	"strconv"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type googleBucketRepo struct {
	bucketName string
	enabled    bool
	creds      []byte
}

// NewGoogleBucketRepo ...
func NewGoogleBucketRepo() domain.GoogleBucketRepository {
	var bucketName = os.Getenv("GOOGLE_BUCKET_NAME")
	enabled, _ := strconv.ParseBool(os.Getenv("GOOGLE_BUCKET_ENABLED"))

	encodedGAC := os.Getenv("GOOGLE_BUCKET_CREDS")
	decodedGAC, _ := base64.StdEncoding.DecodeString(encodedGAC)
	return &googleBucketRepo{
		bucketName: bucketName,
		enabled:    enabled,
		creds:      decodedGAC,
	}
}

func (r *googleBucketRepo) IsEnabled() bool {
	return r.enabled
}

func (r *googleBucketRepo) SupportedType() []string {
	return []string{
		"application/pdf",                                                                               // pdf
		"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", // doc & docx
		"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp", // images
		"video/mp4", "video/mpeg", "video/ogg", "video/webm", // videos
		"audio/mpeg", "audio/ogg", "audio/wav", "audio/webm", // auidos
	}
}

// UploadFile uploads an object.
func (r *googleBucketRepo) UploadFile(ctx context.Context, f multipart.File, objectName string) (*storage.ObjectAttrs, error) {

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(r.creds))
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Upload an object with storage.Writer.
	wc := client.Bucket(r.bucketName).Object(objectName).NewWriter(context.Background())
	if _, err = io.Copy(wc, f); err != nil {
		return nil, fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("Writer.Close: %v", err)
	}

	// set public url
	acl := client.Bucket(r.bucketName).Object(objectName).ACL()
	if err := acl.Set(context.Background(), storage.AllUsers, storage.RoleReader); err != nil {
		return nil, fmt.Errorf("ACLHandle.Set: %v", err)
	}

	return wc.Attrs(), nil
}

func (r *googleBucketRepo) GeneratePublicURL(bucketName, objectName string) string {
	return "https://storage.googleapis.com/" + bucketName + "/" + objectName
}
