package validator

type CreateCategoryPayload struct {
	Name string `json:"name" binding:"required,min=4"`
	Description *string `json:"description" binding:"omitempty,min=10"`
	ParentId string `json:"parentId" binding:"omitempty,uuid"`
}