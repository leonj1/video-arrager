package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"video-editor/app"
)

type MainLayout struct {
	Container   *fyne.Container
	VideoList   *VideoList
	PreviewPane *PreviewPane
	StatusBar   *widget.Label
}

func NewMainLayout(state *app.State, handlers *app.Handlers) *MainLayout {
	videoList := NewVideoList(state)

	previewPane := NewPreviewPane(func(path string) {
		app.PlayVideo(path)
	})

	toolbar := NewToolbar(ToolbarHandlers{
		OnNew:      handlers.OnNew,
		OnAdd:      handlers.OnAddVideos,
		OnRemove:   handlers.OnRemove,
		OnMoveUp:   handlers.OnMoveUp,
		OnMoveDown: handlers.OnMoveDown,
		OnClear:    handlers.OnClear,
		OnExport:   handlers.OnExport,
		OnSave:     handlers.OnSave,
		OnLoad:     handlers.OnLoad,
	})

	header := widget.NewLabel("Video Files (drag to reorder)")
	header.TextStyle = fyne.TextStyle{Bold: true}

	statusBar := widget.NewLabel("No videos")
	statusBar.Alignment = fyne.TextAlignCenter

	listWithHeader := container.NewBorder(
		header,
		nil, nil, nil,
		videoList,
	)

	split := container.NewHSplit(listWithHeader, previewPane)
	split.SetOffset(0.65)

	content := container.NewBorder(
		toolbar,
		statusBar,
		nil, nil,
		split,
	)

	return &MainLayout{
		Container:   content,
		VideoList:   videoList,
		PreviewPane: previewPane,
		StatusBar:   statusBar,
	}
}
