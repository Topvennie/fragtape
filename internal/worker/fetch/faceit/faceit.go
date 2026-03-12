// Package faceit communicates with faceit
package faceit

import (
	"errors"

	"github.com/topvennie/fragtape/pkg/config"
)

var (
	F      *faceit
	apiURL = "https://open.faceit.com/data/v4/"
)

type faceit struct {
	webAPIKey string
}

func Init() error {
	webAPIKey := config.GetDefaultString("worker.fetcher.faceit.api_key", "")
	if webAPIKey == "" {
		return errors.New("no web api key set")
	}

	F = &faceit{
		webAPIKey: webAPIKey,
	}

	return nil
}
