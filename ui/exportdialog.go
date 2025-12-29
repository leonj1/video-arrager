package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"video-editor/app"
)

type ExportDialog struct {
	window            fyne.Window
	transitionSelect  *widget.Select
	durationEntry     *widget.Entry
	onExport          func(options app.ExportOptions)
}

func NewExportDialog(window fyne.Window, onExport func(options app.ExportOptions)) *ExportDialog {
	d := &ExportDialog{
		window:   window,
		onExport: onExport,
	}

	d.transitionSelect = widget.NewSelect([]string{"None", "Fade", "Crossfade"}, nil)
	d.transitionSelect.SetSelected("None")

	d.durationEntry = widget.NewEntry()
	d.durationEntry.SetText("1.0")
	d.durationEntry.SetPlaceHolder("Duration in seconds")

	return d
}

func (d *ExportDialog) Show() {
	transitionLabel := widget.NewLabel("Transition:")
	durationLabel := widget.NewLabel("Duration (seconds):")

	form := container.NewVBox(
		container.NewHBox(transitionLabel, d.transitionSelect),
		container.NewHBox(durationLabel, d.durationEntry),
	)

	dialog.ShowCustomConfirm("Export Options", "Export", "Cancel", form, func(confirmed bool) {
		if !confirmed {
			return
		}

		options := app.ExportOptions{}

		switch d.transitionSelect.Selected {
		case "Fade":
			options.Transition = app.TransitionFade
		case "Crossfade":
			options.Transition = app.TransitionCrossfade
		default:
			options.Transition = app.TransitionNone
		}

		if dur, err := strconv.ParseFloat(d.durationEntry.Text, 64); err == nil && dur > 0 {
			options.TransitionDuration = dur
		} else {
			options.TransitionDuration = 1.0
		}

		if d.onExport != nil {
			d.onExport(options)
		}
	}, d.window)
}
