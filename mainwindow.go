package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Line struct {
	parent uint64
	child  uint64
	line   *widgets.QGraphicsLineItem
}

var links map[uint64]Line

var view *widgets.QGraphicsView

func GetGroupUID(group *widgets.QGraphicsItemGroup) uint64 {
	return group.Data(0).ToULongLong(nil)
}

func GetGroupFromUID(id uint64) *widgets.QGraphicsItemGroup {
	for _, item := range view.Items() {
		if group := item.Group(); group != nil && GetGroupUID(group) == id {
			return group
		}
	}
	return nil
}

func AddLink(parent, child *widgets.QGraphicsItemGroup) *widgets.QGraphicsLineItem {
	// Check if map needs to be created
	if links == nil {
		links = make(map[uint64]Line)
	}
	// Get from (parent) and to (child)
	fromPos := parent.Pos()
	toPos := child.Pos()
	// Create graphics line
	line := widgets.NewQGraphicsLineItem3(
		fromPos.X()+32, fromPos.Y()+32,
		toPos.X()+32, toPos.Y()+32,
		nil,
	)
	// Set the color of it
	line.SetPen(gui.NewQPen3(gui.NewQColor3(0, 255, 0, 255)))
	// Create line data
	parentID := GetGroupUID(parent)
	childID := GetGroupUID(child)
	lineData := Line{
		parentID, childID, line,
	}
	links[parentID] = lineData
	links[childID] = lineData
	// Return the graphics line to add to scene
	return line
}

func GetRandomItemUID() uint64 {
	// TODO: This should guarantee unique, for now, just return random uint64
	return rand.Uint64()
}

func UpdateLinkPos(item *widgets.QGraphicsItemGroup, x, y float64) {
	// Get link
	itemID := GetGroupUID(item)
	link := links[itemID]
	// Error checking
	if link == (Line{}) {
		return
	}
	// If the item is the parent
	isParent := link.parent == itemID
	// Update position of either parent or child
	if isParent {
		pos := GetGroupFromUID(link.child).Pos()
		link.line.SetLine2(x+32, y+32, pos.X()+32, pos.Y()+32)
	} else {
		pos := GetGroupFromUID(link.parent).Pos()
		link.line.SetLine2(pos.X()+32, pos.Y()+32, x+32, y+32)
	}
}

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
	// Create scene and view
	scene := widgets.NewQGraphicsScene(nil)
	view = widgets.NewQGraphicsView2(scene, nil)

	// Set view as central widget
	window.SetCentralWidget(view)

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
	linkRadio := widgets.NewQRadioButton2("Link", nil)
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
		if item != nil && item.Group() != nil {
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
			// Update link if needed
			if _, ok := links[GetGroupUID(movingItem)]; ok {
				UpdateLinkPos(movingItem, movingItem.Pos().X(), movingItem.Pos().Y())
			}
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
			toPos := view.ItemAt(event.Pos()).Group().Pos()
			if toPos.X() == 0 && toPos.Y() == 0 {
				return
			}
			scene.AddItem(AddLink(linkStart, view.ItemAt(event.Pos()).Group()))
			linkStart = nil
		}
	})
	// Show the view
	view.Show()
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

//AddGraphicsItem
func AddGraphicsItem(view *widgets.QGraphicsView, text string, x, y, width, height float64) *widgets.QGraphicsItemGroup {
	group := widgets.NewQGraphicsItemGroup(nil)
	textItem := widgets.NewQGraphicsTextItem2(text, nil)
	shapeItem := widgets.NewQGraphicsRectItem3(0, 0, width, height, nil)
	group.AddToGroup(textItem)
	group.AddToGroup(shapeItem)
	group.SetPos2(x, y)
	group.SetData(0, core.NewQVariant1(GetRandomItemUID()))
	return group
}
