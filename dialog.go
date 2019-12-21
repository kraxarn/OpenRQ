package main

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Dialog is a simpler message box
type Dialog widgets.QMessageBox

// NewDialog creates a new message box without showing it
func NewDialog(parent widgets.QWidget_ITF, title, message string) *Dialog {
	dialog := widgets.NewQMessageBox(parent)
	dialog.SetWindowTitle(title)
	dialog.SetText(message)
	return (*Dialog)(dialog)
}

// SetIcon sets message icon from system theme
func (dialog *Dialog) SetIcon(name string) {
	(*widgets.QMessageBox)(dialog).SetIconPixmap(gui.QIcon_FromTheme(name).Pixmap2(32, 32, 0, 0))
}

// ShowInformation sets the icon as information, button as "OK" and shows it
func (dialog *Dialog) ShowInformation() {
	dialog.SetIcon("messagebox_info")
	dialog.Show()
}
