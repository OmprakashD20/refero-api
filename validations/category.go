package validator

type CreateCategoryPayload struct {
	Name        string  `json:"name" binding:"required,min=4"`
	Description *string `json:"description" binding:"omitempty,min=10"`
	ParentId    string  `json:"parentId" binding:"omitempty,uuid"`
}

type UpdateCategoryPayload struct {
	Name        *string `json:"name" binding:"omitempty,min=4"`
	Description *string `json:"description" binding:"omitempty,min=10"`
	ParentId    *string `json:"parentId" binding:"omitempty,uuid"`
}

type CategoryParam struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type (
	GetCategoryByIDParam    = CategoryParam
	UpdateCategoryByIDParam = CategoryParam
	DeleteCategoryByIDParam = CategoryParam
)
