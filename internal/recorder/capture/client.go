// Package capture converts a highlight to an actual video
package capture

import (
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

type Capturer struct {
	repo      repository.Repository
	highlight repository.Highlight

	dummy bool
}

func New(repo repository.Repository) *Capturer {
	return &Capturer{
		repo:      repo,
		highlight: *repo.NewHighlight(),
		dummy:     config.GetDefaultBool("recorder.dummy_data", false),
	}
}
