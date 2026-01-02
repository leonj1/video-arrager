package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var videoExtensions = []string{".mp4", ".mov", ".avi", ".mkv", ".webm", ".m4v", ".wmv", ".flv"}

func isVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range videoExtensions {
		if ext == e {
			return true
		}
	}
	return false
}

type Handlers struct {
	state  *State
	window fyne.Window
}

func NewHandlers(state *State, window fyne.Window) *Handlers {
	return &Handlers{
		state:  state,
		window: window,
	}
}

func (h *Handlers) OnNew() {
	if h.state.Count() == 0 {
		return
	}

	dialog.ShowConfirm("New Project", "Start a new project? All unsaved changes will be lost.", func(ok bool) {
		if ok {
			h.state.Clear()
		}
	}, h.window)
}

func (h *Handlers) OnAddVideos() {
	log.Println("Opening file dialog...")
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		log.Println("File dialog callback triggered")
		if err != nil {
			log.Printf("File dialog error: %v", err)
			return
		}
		if reader == nil {
			log.Println("File dialog cancelled")
			return
		}

		log.Println("Getting URI path...")
		uri := reader.URI()
		log.Printf("URI: %v", uri)
		path := uri.Path()
		log.Printf("Path: %s", path)
		
		log.Println("Closing reader...")
		reader.Close()
		log.Println("Reader closed")

		if isVideoFile(path) {
			log.Printf("Adding video: %s", path)
			h.AddVideosWithProgress([]string{path})
			log.Println("Video added")
		}
	}, h.window)

	fd.SetFilter(&videoFilter{})
	log.Println("Showing file dialog...")
	fd.Show()
	log.Println("File dialog shown")
}

func (h *Handlers) OnAddFolder() {
	log.Println("Opening folder dialog...")
	fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			log.Printf("Folder dialog error: %v", err)
			return
		}
		if uri == nil {
			log.Println("Folder dialog cancelled")
			return
		}

		folderPath := uri.Path()
		log.Printf("Selected folder: %s", folderPath)

		var videoPaths []string
		err = filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && isVideoFile(path) {
				videoPaths = append(videoPaths, path)
			}
			return nil
		})
		if err != nil {
			log.Printf("Error walking folder: %v", err)
			return
		}

		log.Printf("Found %d video files", len(videoPaths))
		h.AddVideosWithProgress(videoPaths)
	}, h.window)

	fd.Show()
}

func (h *Handlers) AddVideosWithProgress(paths []string) {
	if len(paths) == 0 {
		return
	}

	progressBar := widget.NewProgressBar()
	progressLabel := widget.NewLabel("Loading videos...")
	content := container.NewVBox(progressLabel, progressBar)

	progressDialog := dialog.NewCustomWithoutButtons("Loading Videos", content, h.window)
	progressDialog.Show()

	go func() {
		for i, path := range paths {
			progress := float64(i+1) / float64(len(paths))
			progressBar.SetValue(progress)
			progressLabel.SetText(fmt.Sprintf("Loading video %d of %d...", i+1, len(paths)))

			log.Printf("Adding video: %s", path)
			h.state.AddVideo(path)
		}
		progressDialog.Hide()
	}()
}

func (h *Handlers) OnRemove() {
	h.state.RemoveSelected()
}

func (h *Handlers) OnMoveUp() {
	h.state.MoveUp()
}

func (h *Handlers) OnMoveDown() {
	h.state.MoveDown()
}

func (h *Handlers) OnClear() {
	if h.state.Count() == 0 {
		return
	}

	dialog.ShowConfirm("Clear All", "Remove all videos from the list?", func(ok bool) {
		if ok {
			h.state.Clear()
		}
	}, h.window)
}

func (h *Handlers) OnExport() {
	if h.state.Count() == 0 {
		dialog.ShowInformation("Export", "No videos to export. Add some videos first.", h.window)
		return
	}

	h.showExportOptions()
}

func (h *Handlers) showExportOptions() {
	transitionSelect := widget.NewSelect([]string{"None", "Fade", "Crossfade"}, nil)
	transitionSelect.SetSelected("None")

	durationEntry := widget.NewEntry()
	durationEntry.SetText("1.0")

	form := widget.NewForm(
		widget.NewFormItem("Transition", transitionSelect),
		widget.NewFormItem("Duration (sec)", durationEntry),
	)

	dialog.ShowCustomConfirm("Export Options", "Next", "Cancel", form, func(confirmed bool) {
		if !confirmed {
			return
		}

		options := ExportOptions{}

		switch transitionSelect.Selected {
		case "Fade":
			options.Transition = TransitionFade
		case "Crossfade":
			options.Transition = TransitionCrossfade
		default:
			options.Transition = TransitionNone
		}

		if dur, err := strconv.ParseFloat(durationEntry.Text, 64); err == nil && dur > 0 {
			options.TransitionDuration = dur
		} else {
			options.TransitionDuration = 1.0
		}

		h.showFileSaveDialog(options)
	}, h.window)
}

func (h *Handlers) showFileSaveDialog(options ExportOptions) {
	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, h.window)
			return
		}
		if writer == nil {
			return
		}

		outputPath := writer.URI().Path()
		writer.Close()

		// Remove the empty file created by Fyne's file save dialog
		// so ffmpeg can write to it without issues
		os.Remove(outputPath)

		videos := h.state.GetVideos()
		progress := make(chan ExportProgress)

		progressDialog := dialog.NewCustom("Exporting", "Cancel", widget.NewLabel("Preparing..."), h.window)
		progressDialog.Show()

		go ExportVideos(videos, outputPath, options, progress)

		go func() {
			for p := range progress {
				if p.Error != nil {
					progressDialog.Hide()
					dialog.ShowError(p.Error, h.window)
					return
				}
				if p.Done {
					progressDialog.Hide()
					dialog.ShowInformation("Export Complete", "Video exported successfully to:\n"+outputPath, h.window)
					return
				}
			}
		}()
	}, h.window)

	fd.SetFileName("combined.mp4")
	fd.Show()
}

func (h *Handlers) OnSave() {
	if h.state.Count() == 0 {
		dialog.ShowInformation("Save Project", "No videos to save. Add some videos first.", h.window)
		return
	}

	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, h.window)
			return
		}
		if writer == nil {
			return
		}

		outputPath := writer.URI().Path()
		writer.Close()

		videos := h.state.GetVideos()
		if err := SaveProject(videos, outputPath); err != nil {
			dialog.ShowError(err, h.window)
			return
		}

		dialog.ShowInformation("Save Project", "Project saved successfully.", h.window)
	}, h.window)

	fd.SetFileName("project.json")
	fd.Show()
}

func (h *Handlers) OnLoad() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, h.window)
			return
		}
		if reader == nil {
			return
		}

		path := reader.URI().Path()
		reader.Close()

		videoPaths, err := LoadProject(path)
		if err != nil {
			dialog.ShowError(err, h.window)
			return
		}

		h.state.Clear()

		for _, videoPath := range videoPaths {
			if err := h.state.AddVideo(videoPath); err != nil {
				log.Printf("Failed to load video %s: %v", videoPath, err)
			}
		}

		dialog.ShowInformation("Load Project", "Project loaded successfully.", h.window)
	}, h.window)

	fd.SetFilter(&jsonFilter{})
	fd.Show()
}

type videoFilter struct{}

func (f *videoFilter) Matches(uri fyne.URI) bool {
	return isVideoFile(uri.Path())
}

func (f *videoFilter) Extensions() []string {
	return videoExtensions
}

type jsonFilter struct{}

func (f *jsonFilter) Matches(uri fyne.URI) bool {
	return strings.ToLower(filepath.Ext(uri.Path())) == ".json"
}

func (f *jsonFilter) Extensions() []string {
	return []string{".json"}
}
