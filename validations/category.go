package validator

type CategoryPayload struct {
	Name        string  `json:"name" binding:"required,min=4"`
	Description *string `json:"description" binding:"omitempty,min=10"`
	ParentId    string `json:"parentId" binding:"omitempty,uuid"`
}

type (
	CreateCategoryPayload = CategoryPayload
	UpdateCategoryPayload = CategoryPayload
)

type CategoryParam struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type (
	GetCategoryByIDParam    = CategoryParam
	UpdateCategoryByIDParam = CategoryParam
	DeleteCategoryByIDParam = CategoryParam
)
