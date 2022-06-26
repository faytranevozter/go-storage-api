package http

import (
	"storage-api/domain"
	"storage-api/media/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type mediaHandler struct {
	MediaUseCase domain.MediaUsecase
}

// NewHandler ..
func NewHandler(g *gin.Engine, bu domain.MediaUsecase, prefix string) {
	var handler mediaHandler
	handler.MediaUseCase = bu

	// custom middleware
	mdl := middleware.InitMiddleware()

	// usecase route
	v1 := g.Group(prefix + "/v1")
	qRoute := v1.Group("/media")
	authRoute := qRoute.Use(mdl.Auth())
	{
		authRoute.GET("/list", handler.ListMedia)
		authRoute.POST("/create", handler.CreateMedia)
		authRoute.GET("/detail/:id", handler.DetailMedia)
		authRoute.PUT("/update/:id", handler.UpdateMedia)
		authRoute.DELETE("/delete/:id", handler.DeleteMedia)
	}
}

// ListMedia handler
func (b mediaHandler) ListMedia(c *gin.Context) {
	ctx := c.Request.Context()

	authData := c.MustGet("JWTDATA").(*domain.JWTPayload)

	resp := b.MediaUseCase.ListMedia(ctx, domain.DefaultPayload{
		Query:    c.Request.URL.Query(),
		AuthData: authData,
	})

	c.JSON(resp.Status, resp)
}

// DetailMedia handler
func (b mediaHandler) DetailMedia(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Param("id")

	authData := c.MustGet("JWTDATA").(*domain.JWTPayload)

	resp := b.MediaUseCase.DetailMedia(ctx, domain.DefaultPayload{
		ID:       ID,
		Query:    c.Request.URL.Query(),
		AuthData: authData,
	})
	c.JSON(resp.Status, resp)
}

// CreateMedia handler
func (b mediaHandler) CreateMedia(c *gin.Context) {
	ctx := c.Request.Context()

	payload := domain.UploadMediaRequest{}
	c.Bind(&payload)

	authData := c.MustGet("JWTDATA").(*domain.JWTPayload)

	resp := b.MediaUseCase.UploadMedia(ctx, domain.DefaultPayload{
		Query:    c.Request.URL.Query(),
		Request:  c.Request,
		Payload:  payload,
		AuthData: authData,
	})
	c.JSON(resp.Status, resp)
}

// UpdateMedia handler
func (b mediaHandler) UpdateMedia(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Param("id")

	payload := domain.UpdateMedia{}
	c.Bind(&payload)

	authData := c.MustGet("JWTDATA").(*domain.JWTPayload)

	resp := b.MediaUseCase.UpdateMedia(ctx, domain.DefaultPayload{
		ID:       ID,
		Query:    c.Request.URL.Query(),
		Request:  c.Request,
		Payload:  payload,
		AuthData: authData,
	})
	c.JSON(resp.Status, resp)
}

// DeleteMedia handler
func (b mediaHandler) DeleteMedia(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Param("id")

	authData := c.MustGet("JWTDATA").(*domain.JWTPayload)

	resp := b.MediaUseCase.DeleteMedia(ctx, domain.DefaultPayload{
		ID:       ID,
		Query:    c.Request.URL.Query(),
		AuthData: authData,
	})

	c.JSON(resp.Status, resp)
}
