package hlae

import (
	"fmt"
	"path/filepath"

	"github.com/topvennie/fragtape/internal/database/model"
)

func (h *Hlae) hlaeExecutable() string {
	return filepath.Join(h.hlaePath, "HLAE.exe")
}

// nolint:unused // This is actually used but golang linter is a bit confused by the windows specific files
func (h *Hlae) hlaeHook() string {
	return filepath.Join(h.hlaePath, "x64", "AfxHookSource2.dll")
}

func (h *Hlae) hlaeScript(demo ...model.Demo) string {
	path := filepath.Join(h.hlaePath, "resources", "AfxHookSource2", "snippets", "fragtape")
	if len(demo) > 0 {
		path = filepath.Join(path, fmt.Sprintf("%d.js", demo[0].ID))
	}

	return path
}

func (h *Hlae) ffmpegExecutable() string {
	return filepath.Join(h.ffmpegPath, "ffmpeg", "bin", "ffmpeg.exe")
}

func (h *Hlae) cs2Executable() string {
	return filepath.Join(h.cs2Path, "game", "bin", "win64", "cs2.exe")
}

func (h *Hlae) cs2Dir() string {
	return filepath.Join(h.cs2Path, "game", "csgo")
}

func (h *Hlae) cs2Video() string {
	return filepath.Join(h.cs2Dir(), "recordings", "fragtape")
}

func (h *Hlae) cs2Demo(demo ...model.Demo) string {
	path := filepath.Join(h.cs2Dir(), "demo", "fragtape")
	if len(demo) > 0 {
		path = filepath.Join(path, fmt.Sprintf("%d.dem", demo[0].ID))
	}

	return path
}

func (h *Hlae) cs2Cfg(demo ...model.Demo) string {
	path := filepath.Join(h.cs2Dir(), "cfg", "fragtape")
	if len(demo) > 0 {
		path = filepath.Join(path, fmt.Sprintf("%d.cfg", demo[0].ID))
	}

	return path
}

func (h *Hlae) cs2Tmp() string {
	return filepath.Join(h.cs2Dir(), "tmp", "fragtape")
}
