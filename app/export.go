package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type TransitionType int

const (
	TransitionNone TransitionType = iota
	TransitionFade
	TransitionCrossfade
)

func (t TransitionType) String() string {
	switch t {
	case TransitionFade:
		return "Fade"
	case TransitionCrossfade:
		return "Crossfade"
	default:
		return "None"
	}
}

type ExportOptions struct {
	Transition         TransitionType
	TransitionDuration float64 // in seconds
}

type ExportProgress struct {
	Status string
	Done   bool
	Error  error
}

func ExportVideos(videos []*Video, outputPath string, options ExportOptions, progress chan<- ExportProgress) {
	defer close(progress)

	if len(videos) == 0 {
		progress <- ExportProgress{Error: fmt.Errorf("no videos to export")}
		return
	}

	progress <- ExportProgress{Status: "Preparing export..."}

	if options.Transition == TransitionNone || len(videos) == 1 {
		exportSimple(videos, outputPath, progress)
	} else {
		exportWithTransitions(videos, outputPath, options, progress)
	}
}

func exportSimple(videos []*Video, outputPath string, progress chan<- ExportProgress) {
	tmpFile, err := os.CreateTemp("", "video-list-*.txt")
	if err != nil {
		progress <- ExportProgress{Error: fmt.Errorf("failed to create temp file: %w", err)}
		return
	}
	defer os.Remove(tmpFile.Name())

	for _, video := range videos {
		escaped := strings.ReplaceAll(video.Path, "'", "'\\''")
		fmt.Fprintf(tmpFile, "file '%s'\n", escaped)
	}
	tmpFile.Close()

	progress <- ExportProgress{Status: "Combining videos..."}

	ext := strings.ToLower(filepath.Ext(outputPath))
	args := []string{
		"-f", "concat",
		"-safe", "0",
		"-i", tmpFile.Name(),
	}

	if ext == ".mp4" || ext == ".mov" || ext == ".m4v" {
		args = append(args, "-c", "copy", "-movflags", "+faststart")
	} else {
		args = append(args, "-c", "copy")
	}

	args = append(args, "-y", outputPath)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		progress <- ExportProgress{Error: fmt.Errorf("ffmpeg error: %w\n%s", err, string(output))}
		return
	}

	progress <- ExportProgress{Status: "Export complete!", Done: true}
}

func exportWithTransitions(videos []*Video, outputPath string, options ExportOptions, progress chan<- ExportProgress) {
	progress <- ExportProgress{Status: "Building transition filters..."}

	duration := options.TransitionDuration
	if duration <= 0 {
		duration = 1.0
	}

	// Build ffmpeg command with xfade filter for crossfade
	// or fade filter for fade in/out
	var args []string

	// Add all input files
	for _, video := range videos {
		args = append(args, "-i", video.Path)
	}

	if options.Transition == TransitionCrossfade {
		args = append(args, buildCrossfadeFilter(videos, duration)...)
	} else if options.Transition == TransitionFade {
		args = append(args, buildFadeFilter(videos, duration)...)
	}

	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext == ".mp4" || ext == ".mov" || ext == ".m4v" {
		args = append(args, "-movflags", "+faststart")
	}

	args = append(args, "-y", outputPath)

	progress <- ExportProgress{Status: "Rendering with transitions..."}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		progress <- ExportProgress{Error: fmt.Errorf("ffmpeg error: %w\n%s", err, string(output))}
		return
	}

	progress <- ExportProgress{Status: "Export complete!", Done: true}
}

func buildCrossfadeFilter(videos []*Video, duration float64) []string {
	n := len(videos)
	if n < 2 {
		return nil
	}

	var filterParts []string
	var audioFilterParts []string

	// Calculate offsets for each transition
	offsets := make([]float64, n-1)
	cumulative := 0.0
	for i := 0; i < n-1; i++ {
		cumulative += videos[i].Duration.Seconds() - duration
		offsets[i] = cumulative
	}

	// Build video xfade chain
	lastVideo := "[0:v]"
	for i := 1; i < n; i++ {
		outputLabel := fmt.Sprintf("[v%d]", i)
		if i == n-1 {
			outputLabel = "[vout]"
		}
		filterParts = append(filterParts,
			fmt.Sprintf("%s[%d:v]xfade=transition=fade:duration=%.2f:offset=%.2f%s",
				lastVideo, i, duration, offsets[i-1], outputLabel))
		lastVideo = outputLabel
	}

	// Build audio crossfade chain
	lastAudio := "[0:a]"
	for i := 1; i < n; i++ {
		outputLabel := fmt.Sprintf("[a%d]", i)
		if i == n-1 {
			outputLabel = "[aout]"
		}
		audioFilterParts = append(audioFilterParts,
			fmt.Sprintf("%s[%d:a]acrossfade=d=%.2f%s",
				lastAudio, i, duration, outputLabel))
		lastAudio = outputLabel
	}

	filter := strings.Join(filterParts, ";") + ";" + strings.Join(audioFilterParts, ";")

	return []string{"-filter_complex", filter, "-map", "[vout]", "-map", "[aout]"}
}

func buildFadeFilter(videos []*Video, duration float64) []string {
	n := len(videos)
	if n < 1 {
		return nil
	}

	var filterParts []string
	var audioFilterParts []string

	// Add fade out at end of each video (except last) and fade in at start (except first)
	for i := 0; i < n; i++ {
		videoDur := videos[i].Duration.Seconds()
		fadeOutStart := videoDur - duration

		var vfilter string
		var afilter string

		if i == 0 {
			// First video: fade out only
			vfilter = fmt.Sprintf("[%d:v]fade=t=out:st=%.2f:d=%.2f[v%d]", i, fadeOutStart, duration, i)
			afilter = fmt.Sprintf("[%d:a]afade=t=out:st=%.2f:d=%.2f[a%d]", i, fadeOutStart, duration, i)
		} else if i == n-1 {
			// Last video: fade in only
			vfilter = fmt.Sprintf("[%d:v]fade=t=in:st=0:d=%.2f[v%d]", i, duration, i)
			afilter = fmt.Sprintf("[%d:a]afade=t=in:st=0:d=%.2f[a%d]", i, duration, i)
		} else {
			// Middle videos: both fade in and fade out
			vfilter = fmt.Sprintf("[%d:v]fade=t=in:st=0:d=%.2f,fade=t=out:st=%.2f:d=%.2f[v%d]", i, duration, fadeOutStart, duration, i)
			afilter = fmt.Sprintf("[%d:a]afade=t=in:st=0:d=%.2f,afade=t=out:st=%.2f:d=%.2f[a%d]", i, duration, fadeOutStart, duration, i)
		}

		filterParts = append(filterParts, vfilter)
		audioFilterParts = append(audioFilterParts, afilter)
	}

	// Concat all faded segments
	var concatInputs string
	for i := 0; i < n; i++ {
		concatInputs += fmt.Sprintf("[v%d][a%d]", i, i)
	}
	concatFilter := fmt.Sprintf("%sconcat=n=%d:v=1:a=1[vout][aout]", concatInputs, n)

	filter := strings.Join(filterParts, ";") + ";" + strings.Join(audioFilterParts, ";") + ";" + concatFilter

	return []string{"-filter_complex", filter, "-map", "[vout]", "-map", "[aout]"}
}
