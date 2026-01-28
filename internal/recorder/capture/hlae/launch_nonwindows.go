//go:build !windows

package hlae

import (
	"context"
	"fmt"
)

func (h *Hlae) launch(_ context.Context, _ string) error {
	return fmt.Errorf("HLAE capture requires Windows")
}
