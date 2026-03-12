// Package faceit communicates with faceit
package faceit

import (
	"context"
	"errors"
	"time"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/config"
)

var (
	F      *faceit
	apiURL = "https://open.faceit.com/data/v4/"
)

type faceit struct {
	timeout time.Time

	webAPIKey string
}

func Init() error {
	webAPIKey := config.GetDefaultString("worker.fetcher.faceit.api_key", "")
	if webAPIKey == "" {
		return errors.New("no web api key set")
	}

	F = &faceit{
		timeout:   time.Time{},
		webAPIKey: webAPIKey,
	}

	return nil
}

func (f *faceit) Fetch(ctx context.Context, user model.User) (model.Demo, bool, error) {
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
