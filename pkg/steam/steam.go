// Package steam communicates with steam with either their api or the steam typescript service
package steam

import (
	"errors"
)

var S *steam

type steam struct {
	steamServiceURL string
	webAPIKey       string
}

func Init(steamServiceURL, webAPIKey string) error {
	if S != nil {
		return nil
	}

	if steamServiceURL == "" {
		return errors.New("no steam service url set")
	}
	if webAPIKey == "" {
		return errors.New("no web api key set")
	}

	S = &steam{
		steamServiceURL: steamServiceURL,
		webAPIKey:       webAPIKey,
	}

	return nil
}
