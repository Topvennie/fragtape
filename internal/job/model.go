// Package job contains shared job related logic
package job

import "github.com/topvennie/fragtape/internal/database/model"

type Highlight struct {
	ID int `json:"id"`
}

func (h Highlight) ToModel() *model.Highlight {
	return &model.Highlight{
		ID: h.ID,
	}
}

type Render struct {
	DemoID     int         `json:"demo_id"`
	Highlights []Highlight `json:"highlights"`
}
