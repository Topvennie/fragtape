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
)

const dummyDir = "./internal/recorder/capture/dummydata/"

func (c *Capturer) Start(ctx context.Context, demo model.Demo) error {
	highlightsAll, err := c.highlight.GetByDemo(ctx, demo.ID)
	if err != nil {
		return err
	}

	// If the recorder part of the pipeline failed then it might have already created some highlights
	highlights := []model.Highlight{}
	for _, h := range highlightsAll {
		if h.FileID == "" {
			highlights = append(highlights, *h)
		}
	}

	if len(highlights) == 0 {
		// No highlights
		return nil
	}

	// Actually render it now
	if c.dummy {
		return c.captureDummy(ctx, highlights)
	}

	return capture()
}

func (c *Capturer) captureDummy(ctx context.Context, highlights []model.Highlight) error {
	dummyVideos, err := os.ReadDir(dummyDir)
	if err != nil {
		return fmt.Errorf("read dummy data dir %w", err)
	}
	if len(dummyVideos) == 0 {
		return errors.New("no dummmy videos")
	}

	// Cleanup is done by the recorder loop
	for _, h := range highlights {
		idx := rand.IntN(len(dummyVideos))
		data, err := os.ReadFile(dummyDir + dummyVideos[idx].Name())
		if err != nil {
			return fmt.Errorf("read file %w", err)
		}

		h.FileID = uuid.NewString()

		if err := storage.S.Set(h.FileID, data, 0); err != nil {
			return fmt.Errorf("store file in storage %w", err)
		}

		if err := c.highlight.Update(ctx, h); err != nil {
			return err
		}
	}

	return nil
}

func capture() error {
	return nil
}
