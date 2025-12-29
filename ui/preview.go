package ui

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"video-editor/app"
)

type PreviewPane struct {
	widget.BaseWidget
	container       *fyne.Container
	thumbnail       *canvas.Image
	nameLabel       *widget.Label
	durationLabel   *widget.Label
	resolutionLabel *widget.Label
	sizeLabel       *widget.Label
	playBtn         *widget.Button
	currentPath     string
	onPlay          func(path string)
}

func NewPreviewPane(onPlay func(path string)) *PreviewPane {
	p := &PreviewPane{
		onPlay: onPlay,
	}

	p.thumbnail = canvas.NewImageFromImage(nil)
	p.thumbnail.SetMinSize(fyne.NewSize(320, 180))
	p.thumbnail.FillMode = canvas.ImageFillContain

	placeholder := canvas.NewRectangle(color.NRGBA{R: 40, G: 40, B: 40, A: 255})
	placeholder.SetMinSize(fyne.NewSize(320, 180))

	p.nameLabel = widget.NewLabel("No video selected")
	p.nameLabel.Wrapping = fyne.TextWrapWord
	p.nameLabel.Alignment = fyne.TextAlignCenter

	p.durationLabel = widget.NewLabel("")
	p.durationLabel.Alignment = fyne.TextAlignCenter
	p.durationLabel.TextStyle = fyne.TextStyle{Bold: true}

	p.resolutionLabel = widget.NewLabel("")
	p.resolutionLabel.Alignment = fyne.TextAlignCenter

	p.sizeLabel = widget.NewLabel("")
	p.sizeLabel.Alignment = fyne.TextAlignCenter

	p.playBtn = widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
		if p.currentPath != "" && p.onPlay != nil {
			p.onPlay(p.currentPath)
		}
	})
	p.playBtn.Disable()

	previewHeader := widget.NewLabel("Preview")
	previewHeader.TextStyle = fyne.TextStyle{Bold: true}

	p.container = container.NewVBox(
		previewHeader,
		container.NewStack(placeholder, p.thumbnail),
		p.nameLabel,
		p.durationLabel,
		p.resolutionLabel,
		p.sizeLabel,
		p.playBtn,
	)

	p.ExtendBaseWidget(p)
	return p
}

func (p *PreviewPane) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.container)
}

func (p *PreviewPane) SetVideo(video *app.Video) {
	if video == nil {
		p.currentPath = ""
		p.thumbnail.Image = nil
		p.thumbnail.Refresh()
		p.nameLabel.SetText("No video selected")
		p.durationLabel.SetText("")
		p.resolutionLabel.SetText("")
		p.sizeLabel.SetText("")
		p.playBtn.Disable()
		return
	}

	p.currentPath = video.Path
	p.nameLabel.SetText(video.Name)
	p.durationLabel.SetText(video.DurationString())
	p.resolutionLabel.SetText(video.ResolutionString())
	p.sizeLabel.SetText(video.SizeString())
	p.playBtn.Enable()

	if video.Thumbnail != nil {
		p.thumbnail.Image = video.Thumbnail
	} else {
		p.thumbnail.Image = image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	p.thumbnail.Refresh()
}
