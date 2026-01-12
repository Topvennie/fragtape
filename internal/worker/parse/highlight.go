package parse

import (
	"errors"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/storage"
)

func getHighlights(demo model.Demo) ([]model.Highlight, error) {
	if demo.FileID == "" {
		return nil, errors.New("demo file deleted")
	}

	_, err := storage.S.Get(demo.FileID)
	if err != nil {
		return nil, fmt.Errorf("get demo file %w", err)
	}

	return nil, nil
}

func submitHighlights(_ []model.Highlight) error {
	return nil
}
