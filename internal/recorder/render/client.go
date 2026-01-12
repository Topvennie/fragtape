// Package render converts a highlight to an actual video
package render

import (
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

type client struct {
	highlight repository.Highlight

	dummy bool
}

var C *client

func Init(repo repository.Repository) {
	C = &client{
		highlight: *repo.NewHighlight(),
		dummy:     config.GetDefaultBool("recorder.dummy_data", false),
	}
}
