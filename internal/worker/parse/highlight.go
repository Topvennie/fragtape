package parse

import (
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/storage"
)

type highlight struct {
	HighlightID int `json:"highlight_id"`
	DemoID      int `json:"demo_id"`
}

func (h highlight) toModel() *model.Highlight {
	return &model.Highlight{}
}

func getHighlights(demo model.Demo) ([]highlight, error) {
	if demo.FileID == "" {
		return nil, errors.New("demo file deleted")
	}

	_, err := storage.S.Get(demo.FileID)
	if err != nil {
		return nil, fmt.Errorf("get demo file %w", err)
	}

	return nil, nil
}

func submitHighlights(_ []highlight) error {
	return nil
}
