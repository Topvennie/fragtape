package hlae

import (
	"path/filepath"
)

func (h *Hlae) hlaeExecutable() string {
	return filepath.Join(h.hlaePath, "HLAE.exe")
}

// nolint:unused // Golang linter is a bit consfused by the windows specific files
func (h *Hlae) hlaeHook() string {
	return filepath.Join(h.hlaePath, "x64", "AfxHookSource2.dll")
}

func (h *Hlae) cs2Executable() string {
	return filepath.Join(h.cs2Path, "game", "bin", "win64", "cs2.exe")
}

func (h *Hlae) cs2Dir() string {
	return filepath.Join(h.cs2Path, "game", "csgo")
}

func (h *Hlae) cs2Video() string {
	return filepath.Join(h.cs2Dir(), "video", "fragtape")
}

func (h *Hlae) cs2Demo() string {
	return filepath.Join(h.cs2Dir(), "demo", "fragtape")
}

func (h *Hlae) cs2Cfg() string {
	return filepath.Join(h.cs2Dir(), "cfg", "fragtape")
}
