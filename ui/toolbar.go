package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ToolbarHandlers struct {
	OnNew      func()
	OnAdd      func()
	OnRemove   func()
	OnMoveUp   func()
	OnMoveDown func()
	OnClear    func()
	OnExport   func()
	OnSave     func()
	OnLoad     func()
}

func NewToolbar(handlers ToolbarHandlers) fyne.CanvasObject {
	newBtn := widget.NewButtonWithIcon("New", theme.DocumentCreateIcon(), handlers.OnNew)
	addBtn := widget.NewButtonWithIcon("Add Videos", theme.ContentAddIcon(), handlers.OnAdd)
	removeBtn := widget.NewButtonWithIcon("Remove", theme.ContentRemoveIcon(), handlers.OnRemove)
	upBtn := widget.NewButtonWithIcon("Move Up", theme.MoveUpIcon(), handlers.OnMoveUp)
	downBtn := widget.NewButtonWithIcon("Move Down", theme.MoveDownIcon(), handlers.OnMoveDown)
	clearBtn := widget.NewButtonWithIcon("Clear All", theme.DeleteIcon(), handlers.OnClear)
	exportBtn := widget.NewButtonWithIcon("Export", theme.DocumentSaveIcon(), handlers.OnExport)
	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), handlers.OnSave)
	loadBtn := widget.NewButtonWithIcon("Load", theme.FolderOpenIcon(), handlers.OnLoad)

	return container.NewHBox(
		newBtn,
		widget.NewSeparator(),
		addBtn,
		removeBtn,
		widget.NewSeparator(),
		upBtn,
		downBtn,
		widget.NewSeparator(),
		clearBtn,
		widget.NewSeparator(),
		saveBtn,
		loadBtn,
		widget.NewSeparator(),
		exportBtn,
	)
}
