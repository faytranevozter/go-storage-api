package domain

// UploadMediaRequest structure for creating a media
type UploadMediaRequest struct {
	Source      string `form:"source"`
	Provider    string `form:"provider"`
	URL         string `form:"url"`
	Name        string `form:"name"`
	Title       string `form:"title"`
	Description string `form:"description"`
}

// UpdateMedia structure for updating a media
type UpdateMedia struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}
