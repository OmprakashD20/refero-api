package links

import (
	"errors"
	"net/http"
	"strings"

	errs "github.com/OmprakashD20/refero-api/errors"
	"github.com/OmprakashD20/refero-api/repository"
	"github.com/OmprakashD20/refero-api/types"
	"github.com/OmprakashD20/refero-api/utils"
	validator "github.com/OmprakashD20/refero-api/validations"

	"github.com/gin-gonic/gin"
)

type LinkService struct {
	store types.LinkStore
	txn   types.TransactionStore
}

func NewService(store types.LinkStore, txn types.TransactionStore) *LinkService {
	return &LinkService{store, txn}
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
	linkID, err := s.store.CheckIfLinkExistsByURL(ctx, link.URL, nil)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}

	// If exists, associate the existing link with new categories
	if linkID != nil {
		err := s.txn.Exec(ctx, func(q *repository.Queries) error {
			existingCategories, err := s.store.GetCategoriesForLink(ctx, *linkID, q)
			if err != nil {
				return errs.InternalServerError(errs.WithCause(err))
			}

			existingCategorySet := make(map[string]struct{}, len(existingCategories))
			for _, categoryID := range existingCategories {
				existingCategorySet[categoryID] = struct{}{}
			}

			var mappings []types.LinkCategoryDTO
			for _, categoryID := range link.CategoryIDs {
				if _, exists := existingCategorySet[categoryID]; !exists {
					mappings = append(mappings, types.LinkCategoryDTO{
						LinkID:     *linkID,
						CategoryID: categoryID,
					})
				}
			}

			if len(mappings) > 0 {
				if err := s.store.AddLinkToCategory(ctx, mappings, q); err != nil {
					return errs.InternalServerError(errs.WithCause(err))
				}
			}

			return nil
		})

		if err != nil {
			c.Error(err)
			return
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
	err = s.txn.Exec(ctx, func(q *repository.Queries) error {
		var err error
		linkID, err = s.store.CreateLink(ctx, link, shortUrl, q)
		if err != nil {
			return errs.InternalServerError(errs.WithError(errs.ErrFailedToCreateLink), errs.WithCause(err))
		}

		// Associate the link with its categories
		var mappings []types.LinkCategoryDTO
		for _, categoryID := range link.CategoryIDs {
			mappings = append(mappings, types.LinkCategoryDTO{
				LinkID:     *linkID,
				CategoryID: categoryID,
			})
		}

		if len(mappings) > 0 {
			if err := s.store.AddLinkToCategory(ctx, mappings, q); err != nil {
				return errs.InternalServerError(errs.WithCause(err))
			}
		}

		return nil
	})

	if err != nil {
		c.Error(err)
		return
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
	data, err := s.store.GetLinkByShortURL(ctx, params.ShortURL, nil)
	if err != nil {
		c.Error(errs.NotFound(errs.ErrLinkNotFound))
		return
	}

	c.Redirect(http.StatusMovedPermanently, data.Url)
}

func (s *LinkService) UpdateLinkByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.UpdateLinkByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	link, ok := validator.GetValidatedData[validator.UpdateLinkPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	err := s.txn.Exec(ctx, func(q *repository.Queries) error {
		// Update the link
		if err := s.store.UpdateLinkByID(ctx, params.ID, link, q); err != nil {
			// If link doesn't exists
			if errors.Is(err, errs.ErrLinkNotFound) {
				return errs.NotFound(errs.ErrLinkNotFound)
			}
			return errs.InternalServerError(errs.WithError(errs.ErrFailedToUpdateLink), errs.WithCause(err))
		}

		// Get the existing categories associated with the link
		existingCategories, err := s.store.GetCategoriesForLink(ctx, params.ID, q)
		if err != nil {
			return errs.InternalServerError(errs.WithError(errs.ErrFailedToUpdateLink), errs.WithCause(err))
		}

		existingCategorySet := make(map[string]struct{}, len(existingCategories))
		for _, categoryID := range existingCategories {
			existingCategorySet[categoryID] = struct{}{}
		}

		newCategorySet := make(map[string]struct{}, len(link.CategoryIDs))
		for _, categoryID := range link.CategoryIDs {
			newCategorySet[categoryID] = struct{}{}
		}

		var categoriesToRemove []types.LinkCategoryDTO
		for _, categoryID := range existingCategories {
			if _, exists := newCategorySet[categoryID]; !exists {
				categoriesToRemove = append(categoriesToRemove, types.LinkCategoryDTO{
					LinkID:     params.ID,
					CategoryID: categoryID,
				})
			}
		}

		var categoriesToAdd []types.LinkCategoryDTO
		for _, categoryID := range link.CategoryIDs {
			if _, exists := existingCategorySet[categoryID]; !exists {
				categoriesToAdd = append(categoriesToAdd, types.LinkCategoryDTO{
					LinkID:     params.ID,
					CategoryID: categoryID,
				})
			}
		}

		if len(categoriesToAdd) > 0 {
			if err := s.store.AddLinkToCategory(ctx, categoriesToAdd, q); err != nil {
				return errs.InternalServerError(errs.WithError(errs.ErrFailedToUpdateLink), errs.WithCause(err))
			}
		}
		if len(categoriesToRemove) > 0 {
			if err := s.store.RemoveLinkToCategory(ctx, categoriesToRemove, q); err != nil {
				return errs.InternalServerError(errs.WithError(errs.ErrFailedToUpdateLink), errs.WithCause(err))
			}
		}

		return nil
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *LinkService) DeleteLinkByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.DeleteLinkByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Delete the link
	if err := s.store.DeleteLinkByID(ctx, params.ID, nil); err != nil {
		// If link doesn't exists
		if errors.Is(err, errs.ErrLinkNotFound) {
			c.Error(errs.NotFound(errs.ErrLinkNotFound))
			return
		}

		c.Error(errs.InternalServerError(errs.WithError(errs.ErrFailedToDeleteLink), errs.WithCause(err)))
		return
	}

	c.JSON(http.StatusOK, nil)
}
