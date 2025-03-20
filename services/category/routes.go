package category

import (
	"errors"
	"net/http"

	errs "github.com/OmprakashD20/refero-api/errors"
	"github.com/OmprakashD20/refero-api/types"
	validator "github.com/OmprakashD20/refero-api/validations"

	"github.com/gin-gonic/gin"
)

type CategoryService struct {
	store types.CategoryStore
}

func NewService(store types.CategoryStore) *CategoryService {
	return &CategoryService{store}
}

func (s *CategoryService) SetupCategoryRoutes(api *gin.RouterGroup) {
	api.POST("/", validator.ValidateBody[validator.CreateCategoryPayload](), s.CreateCategoryHandler)

	api.GET("/", s.GetCategoriesHandler)
	api.GET("/:id", validator.ValidateParams[validator.GetCategoryByIDParam](), s.GetCategoryByIDHandler)
	api.GET("/:id/links", validator.ValidateParams[validator.GetLinksForCategoryParams](), s.GetLinksForCategoryHandler)

	api.PUT("/:id", validator.ValidateParams[validator.UpdateCategoryByIDParam](), validator.ValidateBody[validator.UpdateCategoryPayload](), s.UpdateCategoryByIDHandler)

	api.DELETE("/:id", validator.ValidateParams[validator.DeleteCategoryByIDParam](), s.DeleteCategoryByIDHandler)
}

func (s *CategoryService) CreateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	category, ok := validator.GetValidatedData[validator.CreateCategoryPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Check if the category exists
	exists, err := s.store.CheckIfCategoryExistsByName(ctx, category.Name)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}
	// If exists, return error
	if exists {
		c.Error(errs.Conflict(errs.ErrCategoryExists))
		return
	}

	// Create the category
	if err := s.store.CreateCategory(ctx, category); err != nil {
		c.Error(errs.InternalServerError(errs.WithError(errs.ErrFailedToCreateCategory), errs.WithCause(err)))
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *CategoryService) GetCategoriesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all categories from the database
	categories, err := s.store.GetAllCategories(ctx)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}

	// No categories found in the database
	if categories == nil {
		c.JSON(http.StatusOK, []types.CategoryDTO{})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (s *CategoryService) GetCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.GetCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Get category by the Params ID from database
	category, err := s.store.GetCategoryByID(ctx, params.Id)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}

	// No category found with the Params ID
	if category == nil {
		c.Error(errs.NotFound(errs.ErrCategoryNotFound))
		return
	}

	c.JSON(http.StatusOK, category)
}

func (s *CategoryService) UpdateCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.UpdateCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	category, ok := validator.GetValidatedData[validator.UpdateCategoryPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Update the category
	if err := s.store.UpdateCategoryByID(ctx, params.Id, category); err != nil {
		// If category doesn't exists
		if errors.Is(err, errs.ErrCategoryNotFound) {
			c.Error(errs.NotFound(errs.ErrCategoryNotFound))
			return
		}

		c.Error(errs.InternalServerError(errs.WithError(errs.ErrFailedToUpdateCategory), errs.WithCause(err)))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *CategoryService) DeleteCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.DeleteCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Delete the category
	if err := s.store.DeleteCategoryByID(ctx, params.Id); err != nil {
		// If category doesn't exists
		if errors.Is(err, errs.ErrCategoryNotFound) {
			c.Error(errs.NotFound(errs.ErrCategoryNotFound))
			return
		}

		c.Error(errs.InternalServerError(errs.WithError(errs.ErrFailedToDeleteCategory), errs.WithCause(err)))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *CategoryService) GetLinksForCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.GetCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.Error(errs.BadRequest(errs.ErrInvalidPayload))
		return
	}

	// Check if category exists
	exists, err := s.store.CheckIfCategoryExistsByID(ctx, params.Id)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}
	if !exists {
		c.Error(errs.NotFound(errs.ErrCategoryNotFound))
		return
	}

	links, err := s.store.GetLinksForCategory(ctx, params.Id)
	if err != nil {
		c.Error(errs.InternalServerError(errs.WithCause(err)))
		return
	}

	if len(links) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, links)
}
