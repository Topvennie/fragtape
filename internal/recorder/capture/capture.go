// Package capture converts a highlight to an actual video
package capture

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/recorder/capture/dummy"
	"github.com/topvennie/fragtape/internal/recorder/capture/hlae"
	"github.com/topvennie/fragtape/pkg/config"
	"github.com/topvennie/fragtape/pkg/utils"
)

type Capturer struct {
	repo      repository.Repository
	highlight repository.Highlight

	dummy bool

	captureDummy *dummy.Dummy
	captureHLAE  *hlae.Hlae
}

func New(repo repository.Repository) (*Capturer, error) {
	capturer := &Capturer{
		repo:      repo,
		highlight: *repo.NewHighlight(),
		dummy:     config.GetDefaultBool("recorder.dummy_data", false),
	}

	if capturer.dummy {
		cDummy, err := dummy.New(repo)
		if err != nil {
			return nil, err
		}
		capturer.captureDummy = cDummy
	} else {
		cHLAE, err := hlae.New(repo)
		if err != nil {
			return nil, err
		}
		capturer.captureHLAE = cHLAE
	}

	return capturer, nil
}

func (c *Capturer) Capture(ctx context.Context, demo model.Demo) error {
	highlights, err := c.highlight.GetByDemoPopulated(ctx, demo.ID)
	if err != nil {
		return err
	}

	// If the recorder part of the pipeline failed then it might have already created some highlights
	highlights = utils.SliceFilter(highlights, func(h *model.Highlight) bool { return h.FileID == "" })

	if len(highlights) == 0 {
		// No highlights
		return nil
	}
	if !utils.SliceAll(highlights, func(h *model.Highlight) bool { return len(h.Segments) > 0 }) {
		return fmt.Errorf("demo %d has a highlight without segments", demo.ID)
	}

	if c.dummy {
		return c.captureDummy.Capture(ctx, utils.SliceDereference(highlights))
	}

	return c.captureHLAE.Capture(ctx, demo, utils.SliceDereference(highlights))
}
