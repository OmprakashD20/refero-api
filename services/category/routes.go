package category

import (
	"net/http"

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
}

func (s *CategoryService) CreateCategoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	category, ok := validator.GetValidatedData[validator.CreateCategoryPayload](c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Data",
		})
		return
	}

	// Check if the category exists
	exists, err := s.store.CheckIfCategoryExists(ctx, category.Name)
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
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
	})
}
