package capture

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
)

const dummyDir = "./dummydata/"

func (c *client) Start(ctx context.Context, demo model.Demo) error {
	highlights, err := c.highlight.GetByDemo(ctx, demo.ID)
	if err != nil {
		return err
	}

	if len(highlights) == 0 {
		// No highlights
		return nil
	}

	// Actually render it now
	if c.dummy {
		return c.captureDummy(ctx, utils.SliceDereference(highlights))
	}

	return capture()
}

func (c *client) captureDummy(ctx context.Context, highlights []model.Highlight) error {
	dummyVideos, err := os.ReadDir(dummyDir)
	if err != nil {
		return fmt.Errorf("read dummy data dir %w", err)
	}
	if len(dummyVideos) == 0 {
		return errors.New("no dummmy videos")
	}

	return withRollback(withRollbackStruct[model.Highlight]{
		items: highlights,
		do: func(h model.Highlight) error {
			idx := rand.IntN(len(dummyVideos))
			data, err := os.ReadFile(dummyDir + dummyVideos[idx].Name())
			if err != nil {
				return fmt.Errorf("read file %w", err)
			}

			h.FileID = uuid.NewString()

			if err := storage.S.Set(h.FileID, data, 0); err != nil {
				return fmt.Errorf("store file in storage %w", err)
			}

			if err := c.highlight.UpdateFile(ctx, h); err != nil {
				return err
			}

			return nil
		},
		revert: func(h model.Highlight) {
			if h.FileID != "" {
				h.FileID = ""
				_ = c.highlight.UpdateFile(ctx, h)
				_ = storage.S.Delete(h.FileID)
			}
		},
	})
}

func capture() error {
	return nil
}
