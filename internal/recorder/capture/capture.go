// Package capture converts a highlight to an actual video
package capture

import (
	"context"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/recorder/capture/dummy"
	"github.com/topvennie/fragtape/internal/recorder/capture/hlae"
	"github.com/topvennie/fragtape/pkg/config"
)

type Capturer struct {
	repo      repository.Repository
	highlight repository.Highlight

	dummy bool

	cDummy *dummy.Dummy
	cHLAE  *hlae.Hlae
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
		capturer.cDummy = cDummy
	} else {
		cHLAE, err := hlae.New(repo)
		if err != nil {
			return nil, err
		}
		capturer.cHLAE = cHLAE
	}

	return capturer, nil
}

func (c *Capturer) Capture(ctx context.Context, demo model.Demo) error {
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

	if c.dummy {
		err = c.cDummy.Capture(ctx, highlights)
	} else {
		err = c.cHLAE.Capture(ctx, demo, highlights)
	}

	if err != nil {
		return err
	}

	return nil
}
