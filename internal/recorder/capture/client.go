// Package capture converts a highlight to an actual video
package capture

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/pkg/config"
)

type Capturer struct {
	repo      repository.Repository
	highlight repository.Highlight

	dummy bool
}

func New(repo repository.Repository) (*Capturer, error) {
	capturer := &Capturer{
		repo:      repo,
		highlight: *repo.NewHighlight(),
		dummy:     config.GetDefaultBool("recorder.dummy_data", false),
	}

	if err := capturer.ensureRecorderDeps(); err != nil {
		return nil, err
	}

	return capturer, nil
}

func (c *Capturer) ensureRecorderDeps() error {
	if c.dummy {
		return nil
	}

	if runtime.GOOS != "windows" {
		return fmt.Errorf("recorder.dummy_data is false but runtime is %s; HLAE/CS2 are expected to run in Windows", runtime.GOOS)
	}

	if _, err := exec.LookPath("hlae.exe"); err != nil {
		return fmt.Errorf("recorder.dummy_data is false but HLAE was not found in PATH %w", err)
	}

	if _, err := exec.LookPath("cs2.exe"); err != nil {
		return fmt.Errorf("recorder.dummy_data is false but CS2 was not found in PATH  %w", err)
	}

	return nil
}
