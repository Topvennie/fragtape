package dto

import "github.com/topvennie/fragtape/internal/database/model"

type Highlight struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func HighlightDTO(h *model.Highlight) Highlight {
	return Highlight{
		ID:    h.ID,
		Title: h.Title,
	}
}
