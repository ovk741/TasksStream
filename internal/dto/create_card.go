package dto

type CreateCardRequest struct {
	Title       string `json: "title"`
	Description string `json:"description"`
}
