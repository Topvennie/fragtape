package fetch

import (
	"context"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
)

type faceitFetcher struct {
	timeout time.Time
}

// Interface compliance
var _ fetcher = (*faceitFetcher)(nil)

func newFaceitFetcher() *faceitFetcher {
	return &faceitFetcher{
		timeout: time.Now(),
	}
}

func (f *faceitFetcher) fetch(ctx context.Context, user model.User) (model.Demo, bool, error) {
	demo := model.Demo{
		Source: model.DemoSourceFaceit,
	}

	if time.Now().Before(f.timeout) {
		// We're still waiting a bit
		return demo, false, nil
	}

	if user.Setting.FaceitID == "" {
		// User doesn't have faceit linked
		return demo, false, nil
	}

	// TODO: implement
	return demo, false, nil
}
