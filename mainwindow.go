package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func NewMainWindow() (*widgets.QApplication, *widgets.QMainWindow) {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// Create window
	window := widgets.NewQMainWindow(nil, 0)
	// Set minimum and initial size
	window.SetMinimumSize2(960, 540)
	window.Resize2(1280, 720)
	// Center window on screen
	window.SetGeometry(widgets.QStyle_AlignedRect(
		core.Qt__LeftToRight,
		core.Qt__AlignCenter,
		window.Size(),
		app.Desktop().AvailableGeometry(0)))
	// Set a window title
	window.SetWindowTitle("OpenRQ")

	return app, window
}

// AddToolBar adds a tool bar to the specified window
func AddToolBar(window *widgets.QMainWindow) {
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
	fileTool.SetIcon(gui.QIcon_FromTheme("document-properties"))
	fileTool.SetMenu(fileMenu)
	fileTool.SetToolButtonStyle(core.Qt__ToolButtonTextBesideIcon)
	fileTool.SetPopupMode(widgets.QToolButton__InstantPopup)
	fileToolBar.AddWidget(fileTool)

	// Add requirement/solution buttons
	fileToolBar.AddAction2(gui.QIcon_FromTheme("add"), "New Requirement")
	fileToolBar.AddAction2(gui.QIcon_FromTheme("add"), "New Solution")
	spacer := widgets.NewQWidget(nil, 0)
	spacer.SetSizePolicy2(widgets.QSizePolicy__Expanding, widgets.QSizePolicy__Expanding)
	fileToolBar.AddWidget(spacer)
	validate := fileToolBar.AddAction("Validation Engine")
	validate.SetCheckable(true)
	validate.SetIcon(gui.QIcon_FromTheme("system-search"))
}

// CreateLayout creates the main layout widgets
func CreateLayout(window *widgets.QMainWindow) {
	scene := widgets.NewQGraphicsScene(nil)
	font := gui.NewQFont()
	font.SetPointSize(18)
	scene.AddText("No Project Loaded", font)
	view := widgets.NewQGraphicsView2(scene, nil)
	window.SetCentralWidget(view)
	view.Show()
}
