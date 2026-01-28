//go:build windows

package hlae

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sys/windows"
)

func (h *Hlae) launch(ctx context.Context, cfgPath string) error {
	hlaeExe := h.hlaeExecutable()
	hlaeHook := h.hlaeHook()
	cs2Exe := h.cs2Executable()

	cs2CmdLine := strings.Join([]string{
		"-steam",
		"-insecure",
		"-novid",
		"+sv_lan 1",
		"-sw",
		"-w 1920",
		"-h 1080",
		"+fps_max 60",
		"-afxDisableSteamStorage",
		"+exec " + cfgPath,
	}, " ")

	args := []string{
		"-customLoader",
		"-noGui",
		"-autoStart",
		"-hookDllPath", hlaeHook,
		"-programPath", cs2Exe,
		"-cmdLine", cs2CmdLine,
	}

	cmd := exec.CommandContext(ctx, hlaeExe, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
	}

	zap.S().Debug(hlaeExe)
	zap.S().Debug(cmd.Args)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start HLAE %w", err)
	}

	zap.S().Debug("Waiting")

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("HLAE/CS2 exited with error %w", err)
		}
		return nil
	case <-ctx.Done():
		_ = terminateProcessTree(cmd.Process)
		select {
		case <-done:
		case <-time.After(3 * time.Second):
		}

		return ctx.Err()
	}
}

func terminateProcessTree(p *os.Process) error {
	if p == nil {
		return nil
	}

	return exec.Command("taskkill", "/PID", fmt.Sprint(p.Pid), "/T", "/F").Run()
}
