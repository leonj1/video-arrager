package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"video-editor/app"
)

type VideoList struct {
	widget.BaseWidget
	state       *app.State
	container   *fyne.Container
	items       []*videoItem
	dragIndex   int
	dropIndex   int
	isDragging  bool
}

type videoItem struct {
	widget.BaseWidget
	list       *VideoList
	index      int
	background *canvas.Rectangle
	img        *canvas.Image
	label      *widget.Label
	container  *fyne.Container
}

func newVideoItem(list *VideoList) *videoItem {
	item := &videoItem{
		list:       list,
		background: canvas.NewRectangle(color.Transparent),
		img:        canvas.NewImageFromImage(nil),
		label:      widget.NewLabel(""),
	}
	item.img.SetMinSize(fyne.NewSize(120, 68))
	item.img.FillMode = canvas.ImageFillContain
	item.container = container.NewHBox(item.img, item.label)
	item.ExtendBaseWidget(item)
	return item
}

func (v *videoItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewStack(v.background, v.container))
}

func (v *videoItem) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (v *videoItem) Tapped(*fyne.PointEvent) {
	v.list.state.SetSelected(v.index)
}

func (v *videoItem) Dragged(e *fyne.DragEvent) {
	if !v.list.isDragging {
		v.list.isDragging = true
		v.list.dragIndex = v.index
		v.list.state.SetSelected(v.index)
	}

	itemHeight := float32(80)
	totalDrag := e.Position.Y - itemHeight/2
	newIndex := v.list.dragIndex + int(totalDrag/itemHeight)

	if newIndex < 0 {
		newIndex = 0
	}
	if newIndex >= v.list.state.Count() {
		newIndex = v.list.state.Count() - 1
	}

	v.list.dropIndex = newIndex
	v.list.highlightDropTarget()
}

func (v *videoItem) DragEnd() {
	if v.list.isDragging && v.list.dragIndex != v.list.dropIndex {
		v.list.state.Move(v.list.dragIndex, v.list.dropIndex)
	}
	v.list.isDragging = false
	v.list.clearHighlight()
}

func (v *videoItem) update(index int, video *app.Video) {
	v.index = index
	if video.Thumbnail != nil {
		v.img.Image = video.Thumbnail
		v.img.Refresh()
	}

	duration := video.DurationString()
	resolution := video.ResolutionString()

	var info string
	if duration != "" && resolution != "" {
		info = fmt.Sprintf("%d. %s [%s] %s (%s)", index+1, video.Name, duration, resolution, video.SizeString())
	} else if duration != "" {
		info = fmt.Sprintf("%d. %s [%s] (%s)", index+1, video.Name, duration, video.SizeString())
	} else if resolution != "" {
		info = fmt.Sprintf("%d. %s %s (%s)", index+1, video.Name, resolution, video.SizeString())
	} else {
		info = fmt.Sprintf("%d. %s (%s)", index+1, video.Name, video.SizeString())
	}
	v.label.SetText(info)
}

func (v *videoItem) setSelected(selected bool) {
	if selected {
		v.background.FillColor = color.NRGBA{R: 100, G: 149, B: 237, A: 100}
	} else {
		v.background.FillColor = color.Transparent
	}
	v.background.Refresh()
}

func (v *videoItem) setDropTarget(isTarget bool) {
	if isTarget {
		v.background.FillColor = color.NRGBA{R: 100, G: 200, B: 100, A: 100}
	} else {
		v.setSelected(v.list.state.GetSelected() == v.index)
	}
	v.background.Refresh()
}

func NewVideoList(state *app.State) *VideoList {
	vl := &VideoList{
		state:     state,
		container: container.NewVBox(),
		items:     make([]*videoItem, 0),
		dragIndex: -1,
		dropIndex: -1,
	}
	vl.ExtendBaseWidget(vl)
	return vl
}

func (vl *VideoList) CreateRenderer() fyne.WidgetRenderer {
	scroll := container.NewVScroll(vl.container)
	return widget.NewSimpleRenderer(scroll)
}

func (vl *VideoList) Refresh() {
	videos := vl.state.GetVideos()
	selected := vl.state.GetSelected()

	for len(vl.items) < len(videos) {
		item := newVideoItem(vl)
		vl.items = append(vl.items, item)
	}

	vl.container.Objects = nil
	for i, video := range videos {
		item := vl.items[i]
		item.update(i, video)
		item.setSelected(i == selected)
		vl.container.Add(item)
	}

	vl.container.Refresh()
	vl.BaseWidget.Refresh()
}

func (vl *VideoList) highlightDropTarget() {
	for i, item := range vl.items {
		if i < vl.state.Count() {
			item.setDropTarget(i == vl.dropIndex)
		}
	}
}

func (vl *VideoList) clearHighlight() {
	selected := vl.state.GetSelected()
	for i, item := range vl.items {
		if i < vl.state.Count() {
			item.setSelected(i == selected)
		}
	}
}
