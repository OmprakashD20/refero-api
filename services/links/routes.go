package links

import (
	"net/http"
	"strings"

	errs "github.com/OmprakashD20/refero-api/errors"
	"github.com/OmprakashD20/refero-api/types"
	"github.com/OmprakashD20/refero-api/utils"
	validator "github.com/OmprakashD20/refero-api/validations"

	"github.com/gin-gonic/gin"
)

type LinkService struct {
	store types.LinkStore
}

func NewService(store types.LinkStore) *LinkService {
	return &LinkService{store}
}

func (s *LinkService) SetupLinkRoutes(api *gin.RouterGroup) {
	api.POST("/", validator.ValidateBody[validator.CreateLinkPayload](), s.CreateLinkHandler)

	api.GET("/r/:shortUrl", validator.ValidateParams[validator.RedirectLinkParams](), s.RedirectURLHandler)

	api.PUT("/:id", validator.ValidateParams[validator.UpdateLinkByIDParam](), validator.ValidateBody[validator.UpdateLinkPayload](), s.UpdateLinkByIDHandler)

	api.DELETE("/:id", validator.ValidateParams[validator.DeleteLinkByIDParam](), s.DeleteLinkByIDHandler)
}

func (s *LinkService) CreateLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	link, ok := validator.GetValidatedData[validator.CreateLinkPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Check if link exists
	linkId, err := s.store.CheckIfLinkExistsByURL(ctx, link.URL)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}

	// If exists, associate the existing link with new categories
	if linkId != nil {
		existingCategories, err := s.store.GetCategoriesForLink(ctx, *linkId)
		if err != nil {
			c.Error(errs.InternalServerError(errs.WithCause(err)))
			return
		}

		existingCategorySet := make(map[string]struct{}, len(existingCategories))
		for _, categoryID := range existingCategories {
			existingCategorySet[categoryID] = struct{}{}
		}

		var mappings []types.LinkCategoryDTO
		for _, categoryID := range link.CategoryIDs {
			if _, exists := existingCategorySet[categoryID]; !exists {
				mappings = append(mappings, types.LinkCategoryDTO{
					LinkID:     *linkId,
					CategoryID: categoryID,
				})
			}
		}

		if len(mappings) > 0 {
			if err := s.store.AddLinkToCategory(ctx, mappings); err != nil {
				c.Error(errs.InternalServerError(errs.WithCause(err)))
				return
			}
		}

		c.JSON(http.StatusCreated, nil)
		return
	}

	// Clean the URL
	if !strings.HasPrefix(link.URL, "http://") && !strings.HasPrefix(link.URL, "https://") {
		link.URL = "https://" + link.URL
	}

	// Generate the ShortURL
	shortUrl := utils.GenerateShortURL(link.URL)

	// Insert the link
	linkId, err = s.store.CreateLink(ctx, link, shortUrl)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithError(errs.ErrFailedToCreateLink),
			errs.WithCause(err)))
		return
	}

	// Associate the link with its categories
	var mappings []types.LinkCategoryDTO
	for _, categoryID := range link.CategoryIDs {
		mappings = append(mappings, types.LinkCategoryDTO{
			LinkID:     *linkId,
			CategoryID: categoryID,
		})
	}

	if len(mappings) > 0 {
		if err := s.store.AddLinkToCategory(ctx, mappings); err != nil {
			c.Error(errs.InternalServerError(errs.WithCause(err)))
			return
		}
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *LinkService) RedirectURLHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.RedirectLinkParams](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Get the original link using the short url
	data, err := s.store.GetLinkByShortURL(ctx, params.ShortURL)
	if err != nil {
		c.Error(errs.NotFound(errs.ErrLinkNotFound))
		return
	}

	c.Redirect(http.StatusMovedPermanently, data.Url)
}

func (s *LinkService) UpdateLinkByIDHandler(c *gin.Context) {}
func (s *LinkService) DeleteLinkByIDHandler(c *gin.Context) {}
