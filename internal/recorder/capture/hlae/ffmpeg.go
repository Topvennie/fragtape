package hlae

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

// convert converts the recorded clips to actual mp4 clips
// After hlae is finished the files end up in this structure
// {h.cs2Video()}/{demo.id}/{player.id}/take_xxxx
// with xxxx counting upwards by 1, starting from 0000
// Each segment triggers a new recording (increasing xxxx)
// The order in which something is recorder is:
// - sort the segments of each highlight by start tick
// - sort the highglights by start tick of the first segment
// Each directory take_xxxx contains
// - audio.wav           // The audio file
// - fragtape/video.mp4  // The actual video file
//
// For each highlight all segments need to combined into a single video
// and the audio added
func (h *Hlae) convert(ctx context.Context, demo model.Demo, highlights []model.Highlight) error {
	zap.S().Info("Converting")
	// Only keep the highlights with segments
	// In theorie it shouldn't be possible to have a highlight without segments
	highlights = utils.SliceFilter(highlights, func(h model.Highlight) bool { return len(h.Segments) > 0 })

	// Sort all highlights and segments
	for _, h := range highlights {
		slices.SortFunc(h.Segments, func(a, b model.HighlightSegment) int { return a.StartTick - b.StartTick })
	}
	slices.SortFunc(highlights, func(a, b model.Highlight) int { return a.Segments[0].StartTick - b.Segments[0].StartTick })

	// player id -> segment index
	segmentIndexes := make(map[int]int)
	videoPath := h.cs2Video()

	for _, hl := range highlights {
		segmentIdx := segmentIndexes[hl.UserID]

		dirPaths := make([]string, 0, len(hl.Segments))
		for range hl.Segments {
			dirPaths = append(dirPaths, filepath.Join(videoPath, strconv.Itoa(demo.ID), strconv.Itoa(hl.UserID), fmt.Sprintf("take%04d", segmentIdx)))
			segmentIdx++
		}

		segmentIndexes[hl.UserID] = segmentIdx

		outputPath := filepath.Join(h.cs2Tmp(), strconv.Itoa(demo.ID), strconv.Itoa(hl.UserID), fmt.Sprintf("highlight_%d.mp4", hl.ID))

		if err := h.convertHighlight(ctx, dirPaths, outputPath); err != nil {
			return fmt.Errorf("convert highlight %d | %w", hl.ID, err)
		}

		if err := h.saveHighlight(ctx, hl, outputPath); err != nil {
			return fmt.Errorf("save highlight %d | %w", hl.ID, err)
		}
	}

	return nil
}

func (h *Hlae) convertHighlight(ctx context.Context, dirPaths []string, outputPath string) error {
	if len(dirPaths) == 0 {
		return nil
	}

	type takeFiles struct {
		video string
		audio string
	}

	takes := make([]takeFiles, 0, len(dirPaths))
	for _, dirPath := range dirPaths {
		videoPath := filepath.Join(dirPath, "fragtape", "video.mp4")
		audioPath := filepath.Join(dirPath, "audio.wav")

		if err := requireFile(videoPath); err != nil {
			return fmt.Errorf("missing video file in %s | %w", dirPath, err)
		}
		if err := requireFile(audioPath); err != nil {
			return fmt.Errorf("missing audio file in %s | %w", dirPath, err)
		}

		takes = append(takes, takeFiles{
			video: videoPath,
			audio: audioPath,
		})
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create output dir %s | %w", outputPath, err)
	}

	workDir, err := os.MkdirTemp("", "fragtape-convert-*")
	if err != nil {
		return fmt.Errorf("create tmp dir %w", err)
	}
	defer func() {
		_ = os.RemoveAll(workDir)
	}()

	videoListPath := filepath.Join(workDir, "videos.txt")
	audioListPath := filepath.Join(workDir, "audios.txt")
	concatVideoPath := filepath.Join(workDir, "video_concat.mp4")
	concatAudioPath := filepath.Join(workDir, "audio_concat.wav")

	if err := writeFFmpegConcatList(videoListPath, utils.SliceMap(takes, func(t takeFiles) string { return t.video })); err != nil {
		return fmt.Errorf("write video concat list %w", err)
	}
	if err := writeFFmpegConcatList(audioListPath, utils.SliceMap(takes, func(t takeFiles) string { return t.audio })); err != nil {
		return fmt.Errorf("write audio concat list %w", err)
	}

	// Concat all video segments.
	zap.S().Info("Combining video")
	if err := h.runFFmpeg(ctx,
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", videoListPath,
		"-c", "copy",
		concatVideoPath,
	); err != nil {
		return fmt.Errorf("concat videos %w", err)
	}

	// Concat all wav files.
	zap.S().Info("Combining audio")
	if err := h.runFFmpeg(ctx,
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", audioListPath,
		"-c", "copy",
		concatAudioPath,
	); err != nil {
		return fmt.Errorf("concat audios %w", err)
	}

	// Mux final video + audio.
	zap.S().Info("Muxing")
	if err := h.runFFmpeg(ctx,
		"-y",
		"-i", concatVideoPath,
		"-i", concatAudioPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-shortest",
		outputPath,
	); err != nil {
		return fmt.Errorf("mux final mp4 %w", err)
	}

	return nil
}

func (h *Hlae) saveHighlight(ctx context.Context, highlight model.Highlight, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %w", err)
	}

	highlight.FileID = uuid.NewString()

	if err := storage.S.Set(highlight.FileID, data, 0); err != nil {
		return fmt.Errorf("store file in storage %w", err)
	}

	if err := h.highlight.Update(ctx, highlight); err != nil {
		return err
	}

	return nil
}

func requireFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory", path)
	}
	return nil
}

func writeFFmpegConcatList(path string, filePaths []string) error {
	var b strings.Builder
	for _, p := range filePaths {
		// Concat demuxer wants forward slashes and single quotes escaped.
		pp := filepath.ToSlash(p)
		pp = strings.ReplaceAll(pp, "'", "'\\''")
		fmt.Fprintf(&b, "file '%s'\n", pp)
	}

	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func (h *Hlae) runFFmpeg(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, h.ffmpegExecutable(), args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Env = os.Environ()
	cmd.Stdout = io.Writer(&stdoutBuf)
	cmd.Stderr = io.Writer(&stderrBuf)

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderrBuf.String())
		if msg != "" {
			return fmt.Errorf("%w | %s", err, msg)
		}
		return err
	}

	return nil
}
