// Package faceit communicates with faceit
package faceit

import (
	"errors"
	"time"
)

var (
	F      *faceit
	apiURL = "https://open.faceit.com/data/v4/"
)

type faceit struct {
	timeout time.Time

	webAPIKey string
}

func Init(webAPIKey string) error {
	if webAPIKey == "" {
		return errors.New("no web api key set")
	}

	F = &faceit{
		timeout:   time.Time{},
		webAPIKey: webAPIKey,
	}

	return nil
}
