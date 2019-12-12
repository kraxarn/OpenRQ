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
	// Create scene and view
	scene := widgets.NewQGraphicsScene(nil)
	view := widgets.NewQGraphicsView2(scene, nil)
	// Use layout
	splitter := widgets.NewQSplitter(nil)
	linkRadio := widgets.NewQRadioButton2("Link", nil)
	splitter.AddWidget(CreateItemCreator(linkRadio))
	splitter.AddWidget(view)
	splitter.AddWidget(CreateValidationEngineLayout())
	window.SetCentralWidget(splitter)
	// Add example item
	view.SetAcceptDrops(true)
	view.SetAlignment(core.Qt__AlignTop | core.Qt__AlignLeft)
	view.ConnectDragMoveEvent(func(event *gui.QDragMoveEvent) {
		if event.Source() != nil {
			event.AcceptProposedAction()
		}
	})
	// ID of item to add next
	var itemID int
	// What item we're currently moving, if any
	var movingItem *widgets.QGraphicsItemGroup
	// Start position of link
	var linkStart *widgets.QGraphicsItemGroup

	itemSize := 64.0
	view.ConnectDropEvent(func(event *gui.QDropEvent) {
		pos := view.MapToScene(event.Pos())
		scene.AddItem(AddGraphicsItem(view, fmt.Sprintf("Item %v", itemID), pos.X()-(itemSize/2.0), pos.Y()-(itemSize/2.0), itemSize, itemSize))
		itemID = itemID + 1
	})

	view.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		item := view.ItemAt(event.Pos())
		// If an item was found
		if item != nil {
			if linkRadio.IsChecked() {
				// We're creating a link
				linkStart = item.Group()

			} else {
				// We're moving an item
				movingItem = item.Group()
				movingItem.SetOpacity(0.6)
			}
		}
	})
	view.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		if movingItem != nil {
			movingItem.SetPos(view.MapToScene5(event.Pos().X()-32, event.Pos().Y()-32))
		}
	})
	view.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		// We released a button while moving an item
		if movingItem != nil {
			// Reset opacity and remove as moving
			movingItem.SetOpacity(1.0)
			movingItem = nil
		}
		// We released while creating a link
		if linkStart != nil {
			// If we try to link to the empty void
			if view.ItemAt(event.Pos()).Group() == nil {
				linkStart = nil
				return
			}
			fromPos := linkStart.Pos()
			toItem := view.ItemAt(event.Pos()).Group()
			toPos := toItem.Pos()
			if toPos.X() == 0 && toPos.Y() == 0 {
				return
			}
			scene.AddLine2(
				fromPos.X()+32, fromPos.Y()+32,
				toPos.X()+32, toPos.Y()+32,
				gui.NewQPen3(gui.NewQColor3(0, 255, 0, 255)))
			linkStart = nil
		}
	})
	// Show the view
	view.Show()
}

func CreateValidationEngineLayout() *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(widgets.NewQLabel2("Validation Engine", nil, core.Qt__Widget), 0, core.Qt__AlignTop)

	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	widget.SetMaximumWidth(250)
	return widget
}

func CreateGroupBox(title string, childAlignment core.Qt__AlignmentFlag, children ...widgets.QWidget_ITF) *widgets.QGroupBox {
	layout := widgets.NewQVBoxLayout()
	for _, child := range children {
		layout.AddWidget(child, 1, childAlignment)
	}
	group := widgets.NewQGroupBox2(title, nil)
	group.SetLayout(layout)
	return group
}

func CreateItemCreator(linkRadio *widgets.QRadioButton) *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	// Requirement/solution selection
	reqRadio := widgets.NewQRadioButton2("Requirement", nil)
	reqRadio.SetChecked(true)
	layout.AddWidget(CreateGroupBox("Item Type", core.Qt__AlignTop,
		reqRadio, widgets.NewQRadioButton2("Solution", nil), linkRadio,
	), 1, core.Qt__AlignTop)

	shapeList := widgets.NewQListWidget(nil)
	shapeList.SetDragEnabled(true)
	shapeList.AddItem("Square")
	layout.AddWidget(CreateGroupBox("Create", 0, shapeList), 0, 0)

	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	widget.SetMaximumWidth(250)

	return widget
}

func AddGraphicsItem(view *widgets.QGraphicsView, text string, x, y, width, height float64) *widgets.QGraphicsItemGroup {
	group := widgets.NewQGraphicsItemGroup(nil)
	textItem := widgets.NewQGraphicsTextItem2(text, nil)
	shapeItem := widgets.NewQGraphicsRectItem3(0, 0, width, height, nil)
	group.AddToGroup(textItem)
	group.AddToGroup(shapeItem)
	group.SetPos2(x, y)
	return group
}