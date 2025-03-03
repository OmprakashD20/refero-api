package validator

type LinkPayload struct {
	Title       string  `json:"title" binding:"required,min=4"`
	URL         string  `json:"url" binding:"required,url"`
	Description *string `json:"description" binding:"required,min=10"`
	CategoryID  string  `json:"categoryId" binding:"required,uuid"`
}

type (
	CreateLinkPayload = LinkPayload
	UpdateLinkPayload = LinkPayload
)

type LinkParams struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type (
	GetLinkByIDParam    = LinkParams
	UpdateLinkByIDParam = LinkParams
	DeleteLinkByIDParam = LinkParams
)
