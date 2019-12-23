package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Direction int8
const (
	DirAbove Direction = 0
	DirRight Direction = 1
	DirBelow Direction = 2
	DirLeft  Direction = 3
)

func (d Direction) Flip() Direction {
	return d + 2 % 4
}

type Line struct {
	parent int64
	child  int64
	line   *widgets.QGraphicsLineItem
}

var links map[int64][]*Line

var view *widgets.QGraphicsView

// Items opened in an edit window
var openItems map[int64]*widgets.QDockWidget

func IsItemOpen(uid int64) bool {
	_, ok := openItems[uid]
	return ok
}

func CloseItem(uid int64) {
	delete(openItems, uid)
}

// SnapToGrid naps the specified position to the grid
func SnapToGrid(pos *core.QPoint) *core.QPoint {
	// 2^5=32
	const gridSize = 5
	return core.NewQPoint2((pos.X()>>gridSize<<gridSize)-64, (pos.Y()>>gridSize<<gridSize)-32)
}

func CreateEditWidgetFromPos(pos core.QPoint_ITF) (*widgets.QDockWidget, bool) {
	// Get UID
	uid := GetGroupUID(view.ItemAt(pos).Group())
	// Check if already opened
	if IsItemOpen(uid) {
		// We probably want to put it in focus here or something
		return nil, false
	}
	// Open item
	// TODO: For now, assume requirement
	editWindow := CreateEditWidget(uid, TypeRequirement)
	editWindow.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		CloseItem(uid)
	})
	// Set item as being opened
	openItems[uid] = editWindow
	// Return new window
	return editWindow, true
}

func CreateView(window *widgets.QMainWindow, linkRadio *widgets.QRadioButton) *widgets.QGraphicsView {
	// Create scene and view
	scene := widgets.NewQGraphicsScene(nil)
	view = widgets.NewQGraphicsView2(scene, nil)

	// Create open items map
	openItems = make(map[int64]*widgets.QDockWidget)

	// Setup drag-and-drop
	view.SetAcceptDrops(true)
	view.SetAlignment(core.Qt__AlignTop | core.Qt__AlignLeft)
	view.ConnectDragMoveEvent(func(event *gui.QDragMoveEvent) {
		if event.Source() != nil {
			event.AcceptProposedAction()
		}
	})
	// What item we're currently moving, if any
	var movingItem *widgets.QGraphicsItemGroup
	// Start position of link
	var linkStart *widgets.QGraphicsItemGroup

	itemSize := 64.0
	view.ConnectDropEvent(func(event *gui.QDropEvent) {
		pos := view.MapToScene(event.Pos())

		// Add item to database
		// For now, we assume all items are requirements
		db := currentProject.GetData()
		defer db.Close()
		uid, err := db.AddEmptyRequirement()
		if err != nil {
			widgets.QMessageBox_Warning(
				window, "Failed to add item", err.Error(),
				widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			return
		}
		gridPos := SnapToGrid(pos.ToPoint())
		scene.AddItem(AddGraphicsItem(
			fmt.Sprintf("%x", uid), float64(gridPos.X()), float64(gridPos.Y()), itemSize * 2, itemSize, uid))
		if len(openItems) <= 0 {
			openItems[uid], _ = CreateEditWidgetFromPos(gridPos)
			window.AddDockWidget(core.Qt__RightDockWidgetArea, openItems[uid])
		}
	})

	view.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		if event.Button() != core.Qt__LeftButton {
			return
		}
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
			movingItem.SetPos(view.MapToScene(SnapToGrid(event.Pos())))
		}
	})
	view.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__RightButton && view.ItemAt(event.Pos()).Group() != nil {
			// When right clicking item, show edit/delete options
			menu := widgets.NewQMenu(nil)
			// Edit option
			editAction := menu.AddAction2(gui.QIcon_FromTheme("document-edit"), "Edit")
			editAction.ConnectTriggered(func(checked bool) {
				if editWidget, ok := CreateEditWidgetFromPos(view.MapToScene(event.Pos()).ToPoint()); ok {
					window.AddDockWidget(core.Qt__RightDockWidgetArea, editWidget)
				}
			})
			// Delete option
			deleteAction := menu.AddAction2(gui.QIcon_FromTheme("delete"), "Delete")
			deleteAction.ConnectTriggered(func(checked bool) {
				// Hakke write ur delete here
			})
			// Show menu at cursor
			menu.Popup(view.MapToGlobal(event.Pos()), nil)
			return
		}

		// We released a button while moving an item
		if movingItem != nil {
			// Update link if needed
			// Error handling is already taken care of in UpdateLinkPos
			UpdateLinkPos(movingItem, movingItem.Pos().X(), movingItem.Pos().Y())
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
	return view
}

func GetGroupUID(group *widgets.QGraphicsItemGroup) int64 {
	return group.Data(0).ToLongLong(nil)
}

func WhereIsChild(parent, child *core.QPointF) Direction {
	// Above (childY < parentY)
	if child.Y() < parent.Y() {
		return DirAbove
	}
	// Below (childY > parentY)
	if child.Y() > parent.Y() {
		return DirBelow
	}
	// Left (childX < parentX)
	if child.X() < parent.X() {
		return DirLeft
	}
	// Right (childX > parentX)
	if child.X() > parent.X() {
		return DirRight
	}
	// By default, assume below
	return DirBelow
}

func GetLinkOffset(parent, child *core.QPointF) (fromX, fromY, toX, toY float64) {
	switch WhereIsChild(parent, child) {
	case DirAbove:
		// Center X, parent Y top, child Y bottom
		return 64, 32, 0, 64
	case DirBelow:
		// Center X, parent Y bottom, child Y top
		return 64, 32, 64, 0
	case DirRight:
		// parent X right, child X left, center Y
		return 128, 0, 32, 32
	case DirLeft:
		// parent X left, child X right, center Y
		return 0, 128, 32, 32
	}
	// Default
	return 0, 0, 0, 0
}

func AddLink(parent, child *widgets.QGraphicsItemGroup) *widgets.QGraphicsLineItem {
	// Check if map needs to be created
	if links == nil {
		links = make(map[int64][]*Line)
	}
	// Get from (parent) and to (child)
	fromPos := parent.Pos()
	toPos := child.Pos()
	// Create graphics line
	line := widgets.NewQGraphicsLineItem3(
		fromPos.X()+64, fromPos.Y()+32,
		toPos.X()+64, toPos.Y()+32,
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
	links[parentID] = append(links[parentID], &lineData)
	links[childID] = append(links[childID], &lineData)
	// Return the graphics line to add to scene
	return line
}

func UpdateLinkPos(item *widgets.QGraphicsItemGroup, x, y float64) {
	// Get link
	itemID := GetGroupUID(item)
	link, ok := links[itemID]
	// Error checking
	if !ok {
		return
	}
	for _, l := range link {
		// If the item is the parent
		isParent := l.parent == itemID
		// Update position of either parent or child
		if isParent {
			pos := l.line.Line().P2()
			l.line.SetLine2(x+64, y+32, pos.X(), pos.Y())
		} else {
			pos := l.line.Line().P1()
			l.line.SetLine2(pos.X(), pos.Y(), x+64, y+32)
		}
	}
}

func AddGraphicsItem(text string, x, y, width, height float64, uid int64) *widgets.QGraphicsItemGroup {
	group := widgets.NewQGraphicsItemGroup(nil)
	textItem := widgets.NewQGraphicsTextItem2(text, nil)
	shapeItem := widgets.NewQGraphicsRectItem3(0, 0, width, height, nil)
	group.AddToGroup(textItem)
	group.AddToGroup(shapeItem)
	group.SetPos2(x, y)
	group.SetData(0, core.NewQVariant1(uid))
	return group
}
