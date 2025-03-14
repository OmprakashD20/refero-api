package links

import (
	"net/http"
	"strings"

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
	api.GET("/")
	api.GET("/:id", validator.ValidateParams[validator.GetLinkByIDParam]())
	api.PUT("/:id", validator.ValidateParams[validator.UpdateLinkByIDParam](), validator.ValidateBody[validator.UpdateLinkPayload]())
	api.DELETE("/:id", validator.ValidateParams[validator.DeleteLinkByIDParam]())
}

func (s *LinkService) CreateLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	link, ok := validator.GetValidatedData[validator.CreateLinkPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Data",
		})
		return
	}

	// Check if link exists
	linkId, err := s.store.CheckIfLinkExistsByURL(ctx, link.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// If exists, associate the existing link with new categories
	if linkId != nil {
		existingCategories, err := s.store.GetCategoriesForLink(ctx, *linkId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		existingCategorySet := make(map[string]struct{})
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
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Link created successfully",
		})
		return
	}

	// Clean the URL
	if !strings.HasPrefix(link.URL, "http://") && !strings.HasPrefix(link.URL, "https://") {
		link.URL = "https://" + link.URL
	}

	// Generate the ShortURL
	shortUrl := utils.GenerateShortURL(link.URL)

	// Insert the link
	linkId, err = s.store.InsertLink(ctx, link, shortUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to insert link",
		})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Link created successfully",
		"url":     shortUrl,
	})
}

func (s *LinkService) RedirectURLHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.RedirectLinkParams](c, validator.ValidatedParamKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL",
		})
		return
	}

	// Get the original link using the short url 
	data, err := s.store.GetLinkByShortURL(ctx, params.ShortURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
		})
	}

	c.Redirect(http.StatusMovedPermanently, data.Url)
}
