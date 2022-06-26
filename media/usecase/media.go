package usecase

import (
	"context"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"storage-api/app/helpers"
	"storage-api/app/validator"
	"storage-api/domain"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

func (u *uMedia) ListMedia(ctx context.Context, options domain.DefaultPayload) domain.BaseResponse {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	query := options.Query
	pageQuery := query.Get("page")
	pageInt, _ := strconv.Atoi(pageQuery)
	limitQuery := query.Get("limit")
	limitInt, _ := strconv.Atoi(limitQuery)

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	page := int64(pageInt)
	limit := int64(limitInt)
	offset := (page - 1) * limit

	optionsRepo := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	if query.Get("type") != "" && helpers.InArrayString(query.Get("type"), domain.AvailableUploadType) {
		optionsRepo["type"] = query.Get("type")
	}

	if query.Get("title") != "" {
		optionsRepo["title"] = query.Get("title")
	}

	if query.Get("source") != "" {
		optionsRepo["source"] = query.Get("source")
	}

	if query.Get("sort") != "" && helpers.InArrayString(query.Get("sort"), domain.AllowedSort) {
		optionsRepo["sort"] = query.Get("sort")
		if query.Get("dir") != "" {
			optionsRepo["dir"] = query.Get("dir")
		}
	}

	mMongoList, total, _ := u.repoMongo.FetchMedia(ctx, optionsRepo)

	mList := make([]interface{}, 0)
	for _, mMongo := range mMongoList {
		mList = append(mList, mMongo.ToMedia())
	}

	listResponse := domain.ListResponse{
		List:  mList,
		Limit: limit,
		Page:  page,
		Total: total,
	}

	return helpers.SuccessResp(listResponse)
}

func (u *uMedia) DetailMedia(ctx context.Context, options domain.DefaultPayload) domain.BaseResponse {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	ID := options.ID.(string)

	mediaMongos, total, err := u.repoMongo.FetchMedia(ctx, map[string]interface{}{
		"single": true,
		"id":     ID,
	})
	if err != nil || total == 0 {
		return helpers.ErrResp(http.StatusBadRequest, "media not found")
	}

	mediaMongo := mediaMongos[0]

	return helpers.SuccessResp(mediaMongo.ToMedia())
}

func (um *uMedia) UploadMedia(ctx context.Context, options domain.DefaultPayload) domain.BaseResponse {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	var err error
	var publicURL, typeDocument string
	var file multipart.File
	var uploadedFile *multipart.FileHeader

	payload := options.Payload.(domain.UploadMediaRequest)
	request := options.Request

	// validation
	errValidation := make(map[string]string)

	if payload.Source == "" {
		errValidation["source"] = validator.RequiredErr("source")
	} else if !helpers.InArrayString(payload.Source, []string{"upload", "url"}) {
		errValidation["source"] = validator.InArrayErr("source")
	}

	if payload.Source == "upload" {
		// check provider
		if payload.Provider == "" {
			errValidation["provider"] = validator.RequiredErr("provider")
			return helpers.ErrRespVal(errValidation, "validation error")
		} else if !helpers.InArrayString(payload.Provider, domain.AVAILABLE_PROVIDER) {
			errValidation["provider"] = validator.InArrayErr("provider")
			return helpers.ErrRespVal(errValidation, "validation error")
		}

		if (payload.Provider == domain.PROVIDER_GOOGLE_BUCKET && !um.repoGBucket.IsEnabled()) || (payload.Provider == domain.PROVIDER_CLOUDINARY && !um.repoCloudinary.IsEnabled()) {
			errValidation["provider"] = "provider is disabled by administrator"
			return helpers.ErrRespVal(errValidation, "validation error")
		}

		file, uploadedFile, err = request.FormFile("file")
		if err != nil {
			errValidation["file"] = validator.RequiredErr("file")
			return helpers.ErrRespVal(errValidation, "validation error")
		}
		defer file.Close()

		typeDocument = uploadedFile.Header.Get("Content-Type")
		if !helpers.InArrayString(typeDocument, domain.AllowedMimeTypes) {
			errValidation["file"] = validator.InArrayErr("file type")
			return helpers.ErrRespVal(errValidation, "validation error")
		}

		if (payload.Provider == domain.PROVIDER_GOOGLE_BUCKET && !helpers.InArrayString(typeDocument, um.repoGBucket.SupportedType())) || (payload.Provider == domain.PROVIDER_CLOUDINARY && !helpers.InArrayString(typeDocument, um.repoCloudinary.SupportedType())) {
			errValidation["file"] = "file type is not supported by provider"
			return helpers.ErrRespVal(errValidation, "validation error")
		}
	}

	if payload.Source == "url" && payload.URL == "" {
		errValidation["url"] = validator.RequiredErr("url")
		return helpers.ErrRespVal(errValidation, "validation error")
	}

	if payload.Source == "url" {
		typeDocument, _ = helpers.ContentTypeByURL(payload.URL)
		if !helpers.InArrayString(typeDocument, domain.AllowedMimeTypes) {
			errValidation["url"] = validator.InArrayErr("file type")
			return helpers.ErrRespVal(errValidation, "validation error")
		}

		publicURL = payload.URL
		typeDocument = helpers.CategoryByMime(typeDocument)
	}

	if payload.Name == "" && payload.Source == "upload" {
		payload.Name = uploadedFile.Filename
	}

	if len(errValidation) > 0 {
		return helpers.ErrRespVal(errValidation, "validation error")
	}

	now := time.Now().UTC()
	// default value
	mediaMongo := domain.MediaMongo{
		Title:       payload.Title,
		Description: payload.Description,
		Type:        typeDocument,
		Provider:    payload.Provider,
		PublicURL:   publicURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if payload.Source == "upload" {
		switch payload.Provider {
		case domain.PROVIDER_GOOGLE_BUCKET:
			// generating clean name
			extName := filepath.Ext(payload.Name)
			baseName := strings.TrimSuffix(payload.Name, extName)
			slugBaseName := slug.Make(baseName)
			cleanName := strconv.Itoa(int(time.Now().UnixNano())) + "-" + slugBaseName + extName

			objName := helpers.CategoryByMime(typeDocument) + "/" + cleanName
			uploadResp, err := um.repoGBucket.UploadFile(ctx, file, objName)
			if err != nil {
				return helpers.ErrResp(http.StatusBadRequest, err.Error())
			}

			mediaMongo.Type = helpers.CategoryByMime(typeDocument)
			mediaMongo.Provider = payload.Provider
			mediaMongo.PublicURL = "https://storage.googleapis.com/" + uploadResp.Bucket + "/" + uploadResp.Name
			mediaMongo.GoogleBucketStorage = domain.GoogleBucket{
				BucketName: uploadResp.Bucket,
				ID:         uploadResp.Bucket + "/" + uploadResp.Name + "/" + strconv.Itoa(int(uploadResp.Generation)),
				Name:       uploadResp.Name,
				Generation: uploadResp.Generation,
			}
			break
		case domain.PROVIDER_CLOUDINARY:
			dir := ""
			objName := strconv.Itoa(int(time.Now().Unix()))
			uploadResp, err := um.repoCloudinary.UploadFile(ctx, file, dir, objName)
			if err != nil {
				return helpers.ErrResp(http.StatusBadRequest, err.Error())
			}

			mediaMongo.Type = helpers.CategoryByMime(typeDocument)
			mediaMongo.Provider = payload.Provider
			mediaMongo.PublicURL = uploadResp.SecureURL
			mediaMongo.CloudinaryStorage = domain.CloudinaryData{
				AssetID:      uploadResp.AssetID,
				PublicID:     uploadResp.PublicID,
				Version:      uploadResp.Version,
				VersionID:    uploadResp.VersionID,
				Signature:    uploadResp.Signature,
				Width:        uploadResp.Width,
				Height:       uploadResp.Height,
				Format:       uploadResp.Format,
				ResourceType: uploadResp.ResourceType,
			}
			break
		}
	}

	// insert into db
	err = um.repoMongo.InsertMedia(ctx, &mediaMongo)
	if err != nil {
		return helpers.ErrResp(http.StatusBadRequest, "failed inserting media")
	}

	return helpers.SuccessResp(mediaMongo.ToMedia())
}

func (u *uMedia) UpdateMedia(ctx context.Context, options domain.DefaultPayload) domain.BaseResponse {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	ID := options.ID.(string)
	payload := options.Payload.(domain.UpdateMedia)

	mediaMongos, total, err := u.repoMongo.FetchMedia(ctx, map[string]interface{}{
		"single": true,
		"id":     ID,
	})
	if err != nil || total == 0 {
		return helpers.ErrResp(http.StatusBadRequest, "media not found")
	}

	mediaMongo := mediaMongos[0]

	// validation
	errValidation := make(map[string]string)

	// this is validation section

	if len(errValidation) > 0 {
		return helpers.ErrRespVal(errValidation, "validation error")
	}

	// update data
	mediaMongo.Title = payload.Title
	mediaMongo.Description = payload.Description
	mediaMongo.UpdatedAt = time.Now()

	// update
	err = u.repoMongo.UpdateMedia(ctx, &mediaMongo)
	if err != nil {
		return helpers.ErrResp(http.StatusBadRequest, "failed updating media")
	}

	return helpers.SuccessResp(mediaMongo.ToMedia())
}

func (u *uMedia) DeleteMedia(ctx context.Context, options domain.DefaultPayload) domain.BaseResponse {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	ID := options.ID.(string)

	mediaMongos, total, err := u.repoMongo.FetchMedia(ctx, map[string]interface{}{
		"single": true,
		"id":     ID,
	})
	if err != nil || total == 0 {
		return helpers.ErrResp(http.StatusBadRequest, "media not found")
	}

	mediaMongo := mediaMongos[0]

	// delete
	now := time.Now()
	mediaMongo.DeletedAt = &now

	u.repoMongo.UpdateMedia(ctx, &mediaMongo)

	return helpers.SuccessResp(map[string]interface{}{})
}
