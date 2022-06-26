package helpers

import (
	"errors"
	"net/http"
)

func ContentTypeByURL(url string) (string, error) {
	if url == "" {
		return "", errors.New("url cannot be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	return resp.Header.Get("Content-Type"), nil
}

func ExistFileByURL(url string) (bool, error) {
	if url == "" {
		return false, errors.New("url cannot be empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, errors.New("page or media not exist")
	}
	return true, nil
}

func CategoryByMime(contentType string) string {
	if InArrayString(contentType, []string{"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"}) {
		return "document"
	}
	if InArrayString(contentType, []string{"image/jpeg", "image/png", "image/gif", "image/svg+xml", "image/webp"}) {
		return "image"
	}
	if InArrayString(contentType, []string{"video/mp4", "video/mpeg", "video/ogg", "video/webm"}) {
		return "video"
	}
	if InArrayString(contentType, []string{"audio/mpeg", "audio/ogg", "audio/wav", "audio/webm"}) {
		return "audio"
	}

	return "uncategorized"
}
