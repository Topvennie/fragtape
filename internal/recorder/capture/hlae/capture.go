package hlae

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/storage"
)

func (h *Hlae) Capture(ctx context.Context, demo model.Demo, highlights []model.Highlight) error {
	data, err := storage.S.Get(demo.FileID)
	if err != nil {
		return fmt.Errorf("failed to get demo file from storage %w", err)
	}

	// Save demo file somewhere accessible to CS2
	demoPath := filepath.Join(h.cs2Demo(), fmt.Sprintf("%d.dem", demo.ID))
	if err := os.WriteFile(demoPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write demo file %w", err)
	}

	// Create cfg
	cfgPath := filepath.Join(h.cs2Cfg(), fmt.Sprintf("%d.cfg", demo.ID))
	if err := h.buildCfg(cfgPath, demoPath); err != nil {
		return fmt.Errorf("failed to build cfg %w", err)
	}

	// Cleanup
	// TODO: Uncomment
	// defer func() {
	// 	_ = os.Remove(demoPath)
	// 	_ = os.Remove(cfgPath)
	// }()

	return h.launch(ctx, cfgPath)
}

func (h *Hlae) buildCfg(cfgPath, demoPath string) error {
	demoRel, _ := filepath.Rel(h.cs2Dir(), demoPath)
	demoRel = filepath.ToSlash(demoRel)

	cfg := fmt.Sprintf(
		`
echo "[Fragtape] loading demo"
playdemo %s
		`,
		demoRel,
	)

	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		return fmt.Errorf("failed to write cfg file %w", err)
	}

	return nil
}
