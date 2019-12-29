package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// OpenSizeDialog opens a simple dialog where the user can choose a new width and height in 32*32 tiles
func OpenSizeDialog(parent widgets.QWidget_ITF, initial *core.QPoint, accepted func(w, h int)) {
}