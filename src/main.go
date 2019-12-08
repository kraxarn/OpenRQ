package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(1280, 720)
	window.SetWindowTitle("OpenRQ")

	// Add menu bar
	CreateLayout(window)
	AddToolBars(window)

	// make the window visible
	window.Show()

	// Create example project
	NewProject("project.orq")

	// start the main Qt event loop
	app.Exec()
}

func AddToolBars(window *widgets.QMainWindow) {
	// Create tool bar
	fileToolBar := window.AddToolBar3("")
	// Show both icons and text (default is icons only)
	fileToolBar.SetToolButtonStyle(core.Qt__ToolButtonTextBesideIcon)
	// Hide area for dragging it around
	fileToolBar.SetMovable(false)

	// Add file menu
	fileTool := widgets.NewQToolButton(fileToolBar)
	fileMenu := widgets.NewQMenu2("", fileTool)
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-new"), "New...")
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-open"), "Open...")
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-save"), "Save")
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-save-as"), "Save As...")
	fileMenu.AddSeparator()
	fileQuit := fileMenu.AddAction2(gui.QIcon_FromTheme("application-exit"), "Quit")
	fileQuit.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Quit))
	fileQuit.ConnectTriggered(func(checked bool) {
		window.Close()
	})
	fileTool.SetText("File")
	fileTool.SetMenu(fileMenu)
	fileTool.SetPopupMode(widgets.QToolButton__InstantPopup)
	fileToolBar.AddWidget(fileTool)

	// Add requirement/solution buttons
	fileToolBar.AddAction2(gui.QIcon_FromTheme("add"), "New Requirement")
	fileToolBar.AddAction2(gui.QIcon_FromTheme("add"), "New Solution")
	spacer := widgets.NewQWidget(nil, 0)
	spacer.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
	fileToolBar.AddWidget(spacer)
	fileToolBar.AddAction("Validation Engine").SetCheckable(true)
}

func CreateLayout(window *widgets.QMainWindow) {
	scene := widgets.NewQGraphicsScene(nil)
	font := gui.NewQFont()
	font.SetPointSize(18)
	scene.AddText("No Project Loaded", font)
	view := widgets.NewQGraphicsView2(scene, nil)
	window.SetCentralWidget(view)
	view.Show()
}
