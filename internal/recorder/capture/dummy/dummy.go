package dummy

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/storage"
)

const dataDir = "./internal/recorder/capture/dummy/dummydata/"

type Dummy struct {
	highlight repository.Highlight
}

func New(repo repository.Repository) (*Dummy, error) {
	d := &Dummy{
		highlight: *repo.NewHighlight(),
	}

	if err := d.ensureData(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Dummy) ensureData() error {
	data, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("read dummy data dir %w", err)
	}
	if len(data) == 0 {
		return errors.New("no dummy data")
	}

	return nil
}

func (d *Dummy) Capture(ctx context.Context, highlights []model.Highlight) error {
	dummyVideos, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("read dummy data dir %w", err)
	}
	if len(dummyVideos) == 0 {
		return errors.New("no dummmy videos")
	}

	// Cleanup is done by the recorder loop
	for _, h := range highlights {
		idx := rand.IntN(len(dummyVideos))
		data, err := os.ReadFile(dataDir + dummyVideos[idx].Name())
		if err != nil {
			return fmt.Errorf("read file %w", err)
		}

		h.FileID = uuid.NewString()

		if err := storage.S.Set(h.FileID, data, 0); err != nil {
			return fmt.Errorf("store file in storage %w", err)
		}

		if err := d.highlight.Update(ctx, h); err != nil {
			return err
		}
	}

	return nil
}
