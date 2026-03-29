//go:build windows

package hlae

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/topvennie/fragtape/internal/database/model"
	"go.uber.org/zap"
)

const (
	cs2ImageName = "cs2.exe"
)

func (h *Hlae) launch(ctx context.Context, demo model.Demo) error {
	zap.S().Info("Launching cs2")
	hlaeExe := h.hlaeExecutable()
	hlaeHook := h.hlaeHook()
	cs2Exe := h.cs2Executable()
	cfgAbs := h.cs2Cfg(demo)

	cfgRel, err := filepath.Rel(filepath.Join(h.cs2Dir(), "cfg"), cfgAbs)
	if err != nil {
		return fmt.Errorf("make cfg path relative %w", err)
	}
	cfgRel = filepath.ToSlash(cfgRel)

	// Ensure no stale instance interferes with launch or waiting.
	if err := killAllProcessesByName(ctx, cs2ImageName); err != nil {
		return fmt.Errorf("kill existing cs2 processes %w", err)
	}

	// & "C:\Program Files (x86)\HLAE\HLAE.exe" -customLoader -noGui -autoStart -hookDllPath "C:\Program Files (x86)\HLAE\x64\AfxHookSource2.dll" -programPath "C:\Program Files (x86)\Steam\SteamApps\common\Counter-Strike Global Offensive\game\bin\win64\cs2.exe" -cmdLine "-steam -insecure +sv_lan 1 -console -window -w 1920 -h 1080 +fps_max 60 -afxDisableSteamStorage +exec fragtape/{demo id}.cfg"

	cs2CmdLine := strings.Join([]string{
		"-steam",
		"-insecure",
		"+sv_lan 1",
		"-console",
		"-window",
		"-w 1920",
		"-h 1080",
		"+fps_max 60",
		"-afxDisableSteamStorage",
		"+exec " + cfgRel,
	}, " ")

	args := []string{
		"-customLoader",
		"-noGui",
		"-autoStart",
		"-hookDllPath", hlaeHook,
		"-programPath", cs2Exe,
		"-cmdLine", cs2CmdLine,
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, hlaeExe, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start hlae %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				return fmt.Errorf("wait for hlae %w | stdout=%q | stderr=%q", err, stdoutBuf.String(), stderrBuf.String())
			}

			return nil

		case <-ctx.Done():
			if cmd.Process != nil {
				_ = killProcessTreeByPID(context.Background(), cmd.Process.Pid)
			}
			<-done
			return ctx.Err()
		}
	}
}

func killAllProcessesByName(ctx context.Context, imageName string) error {
	pids, err := findAllProcessesByName(ctx, imageName)
	if err != nil {
		return err
	}

	for _, pid := range pids {
		zap.S().Warnf("killing existing process image %s | pid %d", imageName, pid)

		if err := killProcessTreeByPID(ctx, pid); err != nil {
			return fmt.Errorf("kill %s pid %d | %w", imageName, pid, err)
		}
	}

	return nil
}

func killProcessTreeByPID(ctx context.Context, pid int) error {
	cmd := exec.CommandContext(ctx, "taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("taskkill pid %d %w | %s", pid, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func findAllProcessesByName(ctx context.Context, name string) ([]int, error) {
	out, err := exec.CommandContext(
		ctx,
		"tasklist",
		"/FO", "CSV",
		"/NH",
		"/FI", "IMAGENAME eq "+name,
	).Output()
	if err != nil {
		return nil, fmt.Errorf("tasklist for %s %w", name, err)
	}

	s := strings.TrimSpace(string(out))
	if s == "" || strings.Contains(s, "No tasks are running") {
		return nil, nil
	}

	r := csv.NewReader(strings.NewReader(s))

	var pids []int
	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse tasklist csv for %s | %w", name, err)
		}
		if len(record) < 2 {
			return nil, fmt.Errorf("unexpected tasklist output row for %s | %q", name, record)
		}

		pid, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("parse pid %q for %s | %w", record[1], name, err)
		}

		pids = append(pids, pid)
	}

	return pids, nil
}
