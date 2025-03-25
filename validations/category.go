package validator

type CategoryPayload struct {
	Name        string  `json:"name" binding:"required,min=4"`
	Description *string `json:"description" binding:"omitempty,min=10"`
	ParentId    string  `json:"parentId" binding:"omitempty,uuid"`
}

type (
	CreateCategoryPayload = CategoryPayload
	UpdateCategoryPayload = CategoryPayload
)

type CategoryParams struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type (
	GetCategoryByIDParam      = CategoryParams
	UpdateCategoryByIDParam   = CategoryParams
	DeleteCategoryByIDParam   = CategoryParams
	GetLinksForCategoryParams = CategoryParams
)
