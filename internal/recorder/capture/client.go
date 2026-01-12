// Package capture converts a highlight to an actual video
package capture

import (
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

type client struct {
	repo      repository.Repository
	highlight repository.Highlight

	dummy bool
}

var C *client

func Init(repo repository.Repository) {
	C = &client{
		repo:      repo,
		highlight: *repo.NewHighlight(),
		dummy:     config.GetDefaultBool("recorder.dummy_data", false),
	}
}
