package main

import (
	"fmt"
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
		gui.QGuiApplication_Screens()[0].AvailableGeometry()))
	// Set a window title
	window.SetWindowTitle("OpenRQ")

	return app, window
}

var dockValidation *widgets.QDockWidget

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
	// Add "new project" option
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-new"), "New...")
	// Add "open project" option
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-open"), "Open...")
	// Add "save project" option
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-save"), "Save")
	// Add "save project as" option
	fileMenu.AddAction2(gui.QIcon_FromTheme("document-save-as"), "Save As...")
	// Seperation for other stuff
	fileMenu.AddSeparator()
	// About button that shows version and license information
	fileAbout := fileMenu.AddAction2(gui.QIcon_FromTheme("help-about"), "About")
	fileAbout.ConnectTriggered(func(checked bool) {
		aboutMessage := "This version was compiled without proper version information.\nNo version info available."
		if len(versionTagName) > 0 && len(versionCommitHash) > 0 {
			aboutMessage = fmt.Sprintf("Version %v, commit %v", versionTagName, versionCommitHash)
		}
		widgets.QMessageBox_About(window, "About OpenRQ", aboutMessage)
	})
	// Quit option that closes everything, sets default quit keybind
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
	// Open validation engine widget when toggling
	validate.ConnectTriggered(func(checked bool) {
		if checked {
			dockValidation.Show()
		} else {
			dockValidation.Hide()
		}
	})
}

// CreateLayout creates the main layout widgets
func CreateLayout(window *widgets.QMainWindow) {
	// Set view as central widget
	linkRadio := widgets.NewQRadioButton2("Link", nil)
	view := CreateView(window, linkRadio)
	window.SetCentralWidget(view)
	view.Show()
	// Create validation engine dock widget
	dockValidation = widgets.NewQDockWidget("Validation Engine", window, 0)
	dockValidation.SetWidget(CreateValidationEngineLayout())
	// Hide by default
	dockValidation.Hide()
	// Hide close button as that's done from the tool bar
	dockValidation.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetFloatable)
	// Add dock to main window
	window.AddDockWidget(core.Qt__RightDockWidgetArea, dockValidation)

	// Create item type dock widget
	dockItemType := widgets.NewQDockWidget("Item Type", window, 0)
	dockItemType.SetWidget(CreateItemTypeCreator(linkRadio))
	// Hide close button as there's no reason to close it
	dockItemType.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetFloatable)
	// Add dock to main window
	window.AddDockWidget(core.Qt__LeftDockWidgetArea, dockItemType)

	// Create item shape dock widget
	dockItemShape := widgets.NewQDockWidget("Shape", window, 0)
	dockItemShape.SetWidget(CreateItemShapeCreator())
	// Hide close button as there's no reason to close it
	dockItemShape.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetFloatable)
	window.AddDockWidget(core.Qt__LeftDockWidgetArea, dockItemShape)
}

func CreateValidationEngineLayout() *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(widgets.NewQLabel2("Nothing to validate", nil, core.Qt__Widget), 0, core.Qt__AlignTop)

	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	widget.SetMaximumWidth(250)
	widget.SetMinimumWidth(150)
	return widget
}

func CreateVBoxWidget(children ...widgets.QWidget_ITF) *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	for _, child := range children {
		layout.AddWidget(child, 1, 0)
	}
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	return widget
}

func LayoutToWidget(vbox *widgets.QVBoxLayout) *widgets.QWidget {
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(vbox)
	widget.SetMaximumWidth(200)
	widget.SetMinimumWidth(150)
	return widget
}

//CreateItemTypeCreator
func CreateItemTypeCreator(linkRadio *widgets.QRadioButton) *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	// Requirement/solution selection
	reqRadio := widgets.NewQRadioButton2("Requirement", nil)
	reqRadio.SetChecked(true)
	layout.AddWidget(CreateVBoxWidget(
		reqRadio,
		widgets.NewQRadioButton2("Solution", nil),
		linkRadio), 1, core.Qt__AlignTop)
	return LayoutToWidget(layout)
}

//CreateItemShapeCreator
func CreateItemShapeCreator() *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	shapeList := widgets.NewQListWidget(nil)
	shapeList.SetDragEnabled(true)
	shapeList.AddItem("Square")
	layout.AddWidget(CreateVBoxWidget(shapeList), 0, 0)
	return LayoutToWidget(layout)
}
