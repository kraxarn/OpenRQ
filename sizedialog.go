package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// OpenSizeDialog opens a simple dialog where the user can choose a new width and height in 32*32 tiles
func OpenSizeDialog(parent widgets.QWidget_ITF, initial *core.QPoint, accepted func(w, h int)) {
	// Create dialog
	dialog := widgets.NewQDialog(parent, 1)
	dialog.SetParent(parent)
	dialog.SetFixedSize2(380, 150)
	dialog.SetWindowTitle("Select Size")

	// Widgets used by other widgets
	sizeLabel := widgets.NewQLabel2(
		fmt.Sprintf("Actual size: %vx%v", initial.X()*32, initial.Y()*32), nil, 0)
	width := widgets.NewQSpinBox(nil)
	height := widgets.NewQSpinBox(nil)

	// Main vertical layout
	layout := widgets.NewQVBoxLayout()
	// Grid for width and height selections
	grid := widgets.NewQGridLayout(nil)
	// Width spin box
	grid.AddWidget(widgets.NewQLabel2("Width:", nil, 0), 0, 0, 0)
	width.SetSuffix(" tiles")
	width.SetMinimum(1)
	width.SetMaximum(100)
	width.SetValue(initial.X())
	width.ConnectValueChanged(func(i int) {
		sizeLabel.SetText(fmt.Sprintf("Actual size: %vx%v", i*32, height.Value()*32))
	})
	grid.AddWidget(width, 0, 1, 0)
	// Height spin box
	grid.AddWidget(widgets.NewQLabel2("Height:", nil, 0), 1, 0, 0)
	height.SetSuffix(" tiles")
	height.SetMinimum(1)
	height.SetMaximum(100)
	height.SetValue(initial.Y())
	height.ConnectValueChanged(func(i int) {
		sizeLabel.SetText(fmt.Sprintf("Actual size: %vx%v", width.Value()*32, i*32))
	})
	grid.AddWidget(height, 1, 1, 0)
	// Add grid layout to main layout
	gridWidget := widgets.NewQWidget(nil, 0)
	gridWidget.SetLayout(grid)
	layout.AddWidget(gridWidget, 0, 0)

	// Buttons for cancel and resize
	buttons := widgets.NewQHBoxLayout()
	buttons.AddWidget(sizeLabel, 1, 0)
	btnCancel := widgets.NewQPushButton2("Cancel", nil)
	btnCancel.ConnectClicked(func(checked bool) {
		// When pressing cancel, just close the dialog
		dialog.Close()
	})
	buttons.AddWidget(btnCancel, 0, 0)
	btnResize := widgets.NewQPushButton2("Resize", nil)
	btnResize.ConnectClicked(func(checked bool) {
		// When pressing resize, close, but also call accepted function
		dialog.Close()
		accepted(width.Value(), height.Value())
	})
	buttons.AddWidget(btnResize, 0, 0)
	buttonsWidget := widgets.NewQWidget(nil, 0)
	buttonsWidget.SetLayout(buttons)
	layout.AddWidget(buttonsWidget, 0, 0)

	// Set dialog layout and show it
	dialog.SetLayout(layout)
	dialog.Open()
}