package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	appPkg "video-arranger/app"
	"video-arranger/ui"
)

var videoExtensions = map[string]bool{
	".mp4": true, ".mov": true, ".avi": true, ".mkv": true,
	".webm": true, ".m4v": true, ".wmv": true, ".flv": true,
}

func isVideoFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return videoExtensions[ext]
}

func main() {
	a := app.NewWithID("com.videoarranger.app")
	window := a.NewWindow("Video Arranger")

	state := appPkg.NewState()
	handlers := appPkg.NewHandlers(state, window)
	layout := ui.NewMainLayout(state, handlers)

	state.SetOnChange(func() {
		layout.VideoList.Refresh()

		selected := state.GetSelected()
		videos := state.GetVideos()
		if selected >= 0 && selected < len(videos) {
			layout.PreviewPane.SetVideo(videos[selected])
		} else {
			layout.PreviewPane.SetVideo(nil)
		}

		// Update status bar
		count := state.Count()
		if count == 0 {
			layout.StatusBar.SetText("No videos")
		} else {
			totalDuration := state.TotalDurationString()
			if totalDuration != "" {
				layout.StatusBar.SetText(fmt.Sprintf("%d videos | Total: %s", count, totalDuration))
			} else {
				layout.StatusBar.SetText(fmt.Sprintf("%d videos", count))
			}
		}
	})

	// Keyboard shortcuts
	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		handlers.OnNew()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		handlers.OnAddVideos()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		handlers.OnSave()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyE,
		Modifier: fyne.KeyModifierSuper,
	}, func(_ fyne.Shortcut) {
		handlers.OnExport()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyUp,
	}, func(_ fyne.Shortcut) {
		handlers.OnMoveUp()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyDown,
	}, func(_ fyne.Shortcut) {
		handlers.OnMoveDown()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyDelete,
	}, func(_ fyne.Shortcut) {
		handlers.OnRemove()
	})

	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName: fyne.KeyBackspace,
	}, func(_ fyne.Shortcut) {
		handlers.OnRemove()
	})

	// Drag and drop support
	window.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		var paths []string
		for _, uri := range uris {
			path := uri.Path()
			if isVideoFile(path) {
				paths = append(paths, path)
			}
		}
		if len(paths) > 0 {
			handlers.AddVideosWithProgress(paths)
		}
	})

	window.SetContent(layout.Container)
	window.Resize(fyne.NewSize(900, 500))
	window.ShowAndRun()
}
