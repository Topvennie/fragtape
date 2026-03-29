//go:build !windows

package hlae

import (
	"context"
	"fmt"

	"github.com/topvennie/fragtape/internal/database/model"
)

func (h *Hlae) launch(_ context.Context, _ model.Demo) error {
	return fmt.Errorf("HLAE capture requires Windows")
}
