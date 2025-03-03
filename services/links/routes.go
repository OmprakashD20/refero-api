package links

import (
	"github.com/OmprakashD20/refero-api/types"
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
	api.POST("/", validator.ValidateBody[validator.CreateLinkPayload]())
	api.GET("/")
	api.GET("/:id", validator.ValidateParams[validator.GetLinkByIDParam]())
	api.PUT("/:id", validator.ValidateParams[validator.UpdateLinkByIDParam](), validator.ValidateBody[validator.UpdateLinkPayload]())
	api.DELETE("/:id", validator.ValidateParams[validator.DeleteLinkByIDParam]())
}
