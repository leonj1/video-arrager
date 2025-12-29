package app

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Video struct {
	Path      string
	Name      string
	Size      int64
	Duration  time.Duration
	Width     int
	Height    int
	Thumbnail image.Image
}

func NewVideo(path string) (*Video, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	video := &Video{
		Path: path,
		Name: filepath.Base(path),
		Size: info.Size(),
	}

	if thumb, err := ExtractThumbnail(path); err == nil {
		video.Thumbnail = thumb
	}

	if duration, err := ExtractDuration(path); err == nil {
		video.Duration = duration
	}

	if width, height, err := ExtractResolution(path); err == nil {
		video.Width = width
		video.Height = height
	}

	return video, nil
}

func ExtractDuration(videoPath string) (time.Duration, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	durationStr := strings.TrimSpace(out.String())
	seconds, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds * float64(time.Second)), nil
}

func ExtractResolution(videoPath string) (int, int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=s=x:p=0",
		videoPath)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, 0, err
	}

	parts := strings.Split(strings.TrimSpace(out.String()), "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid resolution format")
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}

func (v *Video) ResolutionString() string {
	if v.Width == 0 || v.Height == 0 {
		return ""
	}
	return fmt.Sprintf("%dx%d", v.Width, v.Height)
}

func (v *Video) DurationString() string {
	if v.Duration == 0 {
		return ""
	}

	total := int(v.Duration.Seconds())
	hours := total / 3600
	minutes := (total % 3600) / 60
	seconds := total % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func (v *Video) SizeString() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case v.Size >= GB:
		return formatSize(float64(v.Size)/GB, "GB")
	case v.Size >= MB:
		return formatSize(float64(v.Size)/MB, "MB")
	case v.Size >= KB:
		return formatSize(float64(v.Size)/KB, "KB")
	default:
		return formatSize(float64(v.Size), "B")
	}
}

func formatSize(size float64, unit string) string {
	if size == float64(int(size)) {
		return fmt.Sprintf("%d %s", int(size), unit)
	}
	return fmt.Sprintf("%.1f %s", size, unit)
}
