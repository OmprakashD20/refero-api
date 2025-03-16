package category

import (
	"errors"
	"net/http"

	db_errors "github.com/OmprakashD20/refero-api/errors"
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
	api.PUT("/:id", validator.ValidateParams[validator.UpdateCategoryByIDParam](), validator.ValidateBody[validator.UpdateCategoryPayload](), s.UpdateCategoryByIDHandler)
	api.DELETE("/:id", validator.ValidateParams[validator.DeleteCategoryByIDParam](), s.DeleteCategoryByIDHandler)
	api.GET("/:id/links", validator.ValidateParams[validator.GetLinksForCategoryParams](), s.GetLinksForCategoryHandler)
}

func (s *CategoryService) CreateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	category, ok := validator.GetValidatedData[validator.CreateCategoryPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Data",
		})
		return
	}

	// Check if the category exists
	exists, err := s.store.CheckIfCategoryExistsByName(ctx, category.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	// If exists, return error
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Category already exists",
		})
		return
	}

	// Create the category
	if err := s.store.CreateCategory(ctx, category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *CategoryService) GetCategoriesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Get all categories from the database
	categories, err := s.store.GetAllCategories(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Category",
		})
		return
	}

	// Get category by the Params ID from database
	category, err := s.store.GetCategoryByID(ctx, params.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	// No category found with the Params ID
	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (s *CategoryService) UpdateCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.UpdateCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Category",
		})
		return
	}

	category, ok := validator.GetValidatedData[validator.UpdateCategoryPayload](c, validator.ValidatedBodyKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Data",
		})
		return
	}

	// Update the category
	if err := s.store.UpdateCategoryByID(ctx, params.Id, category); err != nil {
		// If category doesn't exists
		if errors.Is(err, db_errors.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update the category",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *CategoryService) DeleteCategoryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.UpdateCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Category",
		})
		return
	}

	// Delete the category
	if err := s.store.DeleteCategoryByID(ctx, params.Id); err != nil {
		// If category doesn't exists
		if errors.Is(err, db_errors.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Category not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete the category",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *CategoryService) GetLinksForCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	params, ok := validator.GetValidatedData[validator.GetCategoryByIDParam](c, validator.ValidatedParamKey)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Category",
		})
		return
	}

	// Check if category exists
	exists, err := s.store.CheckIfCategoryExistsByID(ctx, params.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category does not exists",
		})
		return
	}

	links, err := s.store.GetLinksForCategory(ctx, params.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	if len(links) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, links)
}
