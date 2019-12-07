package main

import (
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(1280, 720)
	window.SetWindowTitle("OpenRQ")

	// Add menu bar
	AddMenuBars(window)
	AddToolBars(window)

	// make the window visible
	window.Show()

	// Create example project
	NewProject("project.orq")

	// start the main Qt event loop
	app.Exec()
}

func AddMenuBars(window *widgets.QMainWindow) {
	// File menu
	fileMenu := window.MenuBar().AddMenu2("File")
	fileMenu.AddAction("New...")
	fileMenu.AddAction("Open...")
	fileMenu.AddAction("Save")
	fileMenu.AddAction("Save As...")
	fileMenu.AddSeparator()
	fileMenu.AddAction("Close").ConnectTriggered(func(checked bool) {
		window.Close()
	})
	// Add menu
	addMenu := window.MenuBar().AddMenu2("Add")
	addMenu.AddAction("Requirement")
	addMenu.AddAction("Solution")
	// View menu
	viewMenu := window.MenuBar().AddMenu2("View")
	viewMenu.AddAction("Validation Engine")
}

func AddToolBars(window *widgets.QMainWindow) {
	fileToolBar := window.AddToolBar3("")
	fileToolBar.SetMovable(false)
	fileToolBar.AddAction("New Requirement")
	fileToolBar.AddAction("New Solution")
	spacer := widgets.NewQWidget(nil, 0)
	spacer.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
	fileToolBar.AddWidget(spacer)
	fileToolBar.AddAction("Validation Engine").SetCheckable(true)
}
