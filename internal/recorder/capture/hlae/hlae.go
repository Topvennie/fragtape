// Package hlae provides functionality to capture videos using HLAE and CS2.
package hlae

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

const (
	tmpDirVideo = "./tmp_recorder/video"
	tmpDirDemo  = "./tmp_recorder/demo"
)

type Hlae struct {
	hlaePath string
	cs2Path  string
}

func New(repo repository.Repository) (*Hlae, error) {
	h := &Hlae{
		hlaePath: config.GetString("recorder.hlae_path"),
		cs2Path:  config.GetString("recorder.cs2_path"),
	}

	if err := h.ensureDeps(); err != nil {
		return nil, err
	}

	if err := h.ensureTmpDirs(); err != nil {
		return nil, err
	}

	if err := h.cleanTmpDir(tmpDirVideo); err != nil {
		return nil, err
	}

	if err := h.cleanTmpDir(tmpDirDemo); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Hlae) ensureDeps() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("runtime is %s; HLAE/CS2 are expected to run in Windows", runtime.GOOS)
	}

	if _, err := exec.LookPath(hlaeExecutable(h.hlaePath)); err != nil {
		return fmt.Errorf("HLAE was not found %w", err)
	}

	if _, err := exec.LookPath(cs2Exectuable(h.cs2Path)); err != nil {
		return fmt.Errorf("CS2 was not found %w", err)
	}

	return nil
}

func (h *Hlae) ensureTmpDirs() error {
	if err := os.MkdirAll(tmpDirVideo, os.ModePerm); err != nil {
		return fmt.Errorf("create tmp video dir %w", err)
	}

	if err := os.MkdirAll(tmpDirDemo, os.ModePerm); err != nil {
		return fmt.Errorf("create tmp demo dir %w", err)
	}

	return nil
}

func (h *Hlae) cleanTmpDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read tmp dir %w", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			if err := os.Remove(fmt.Sprintf("%s/%s", dir, e.Name())); err != nil {
				return fmt.Errorf("remove tmp subdir %w", err)
			}

			continue
		}

		if err := os.Remove(fmt.Sprintf("%s/%s", dir, e.Name())); err != nil {
			return fmt.Errorf("remove tmp file %w", err)
		}
	}

	return nil
}

func (h *Hlae) Capture(ctx context.Context, demo model.Demo, highlights []model.Highlight) error {
	return nil
}
