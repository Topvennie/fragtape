// Package hlae provides functionality to capture videos using HLAE and CS2.
package hlae

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

type Hlae struct {
	hlaePath   string
	ffmpegPath string
	cs2Path    string

	highlight repository.Highlight
	user      repository.User
}

func New(repo repository.Repository) (*Hlae, error) {
	h := &Hlae{
		hlaePath:   config.GetString("recorder.hlae_path"),
		ffmpegPath: config.GetString("recorder.ffmpeg_path"),
		cs2Path:    config.GetString("recorder.cs2_path"),
		highlight:  *repo.NewHighlight(),
		user:       *repo.NewUser(),
	}

	if err := h.ensureDeps(); err != nil {
		return nil, err
	}

	if err := h.ensureTmpDirs(); err != nil {
		return nil, err
	}

	if err := h.cleanTmpDir(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Hlae) ensureDeps() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("runtime is %s, HLAE/CS2 are expected to run in Windows", runtime.GOOS)
	}

	if _, err := os.Stat(h.hlaeExecutable()); err != nil {
		return fmt.Errorf("HLAE was not found %w", err)
	}

	if _, err := os.Stat(h.hlaeHook()); err != nil {
		return fmt.Errorf("HLAE hook was not found %w", err)
	}

	if _, err := os.Stat(h.ffmpegExecutable()); err != nil {
		return fmt.Errorf("FFMPEG was not found %w", err)
	}

	if _, err := os.Stat(h.cs2Executable()); err != nil {
		return fmt.Errorf("CS2 was not found %w", err)
	}

	return nil
}

func (h *Hlae) ensureTmpDirs() error {
	if err := os.MkdirAll(h.hlaeScript(), 0o755); err != nil {
		return fmt.Errorf("create cfg dir %w", err)
	}

	if err := os.MkdirAll(h.cs2Video(), 0o755); err != nil {
		return fmt.Errorf("create video dir %w", err)
	}

	if err := os.MkdirAll(h.cs2Demo(), 0o755); err != nil {
		return fmt.Errorf("create demo dir %w", err)
	}

	if err := os.MkdirAll(h.cs2Cfg(), 0o755); err != nil {
		return fmt.Errorf("create cfg dir %w", err)
	}

	if err := os.MkdirAll(h.cs2Tmp(), 0o755); err != nil {
		return fmt.Errorf("create tmp dir %w", err)
	}

	return nil
}

func (h *Hlae) cleanTmpDir() error {
	clean := func(dir string) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return fmt.Errorf("read tmp dir %w", err)
		}

		for _, e := range entries {
			p := filepath.Join(dir, e.Name())
			if err := os.RemoveAll(p); err != nil {
				return fmt.Errorf("remove %q: %w", p, err)
			}
		}

		return nil
	}

	if err := clean(h.hlaeScript()); err != nil {
		return fmt.Errorf("clean script dir %w", err)
	}

	if err := clean(h.cs2Video()); err != nil {
		return fmt.Errorf("clean video dir %w", err)
	}

	if err := clean(h.cs2Demo()); err != nil {
		return fmt.Errorf("clean demo dir %w", err)
	}

	if err := clean(h.cs2Cfg()); err != nil {
		return fmt.Errorf("clean cfg dir %w", err)
	}

	if err := clean(h.cs2Tmp()); err != nil {
		return fmt.Errorf("clean tmp dir %w", err)
	}

	return nil
}
