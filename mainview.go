package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Line struct {
	parent int64
	child  int64
	line   *widgets.QGraphicsLineItem
	dir    *widgets.QGraphicsPolygonItem
}

var links map[int64][]*Line

var view *widgets.QGraphicsView

var backgroundColor *gui.QColor

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
	scenePos := view.MapToScene(pos).ToPoint()
	return view.MapFromScene(core.NewQPointF3(
		float64((scenePos.X()>>gridSize<<gridSize)-64), float64((scenePos.Y()>>gridSize<<gridSize)-32)))
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

func CreateView(window *widgets.QMainWindow, linkBtn *widgets.QToolButton) *widgets.QGraphicsView {
	// Create scene and view
	scene := widgets.NewQGraphicsScene(nil)
	view = widgets.NewQGraphicsView2(scene, nil)
	// Get default window background color
	backgroundColor = window.Palette().Color2(window.BackgroundRole())
	// Create open items map
	openItems = make(map[int64]*widgets.QDockWidget)

	// Load items from database
	{
		db := currentProject.Data()
		defer db.Close()
		items, err := db.Items()
		if err != nil {
			fmt.Println("error: failed to get saved items:", err)
		} else {
			var x, y, w, h int
			for _, item := range items {
				x, y = item.Pos()
				w, h = item.Size()
				scene.AddItem(NewGraphicsItem(
					fmt.Sprintf("%x\n%v", item.ID(), item.Description()), x, y, w, h, item.ID()))
			}
		}
		// Get links
		links, err := db.Links()
		if err != nil {
			fmt.Println("error: failed to get saved links:", err)
		} else {
			for parent, child := range links {
				// Find parent and child
				var parentItem, childItem *widgets.QGraphicsItemGroup
				for _, item := range view.Items() {
					group := item.Group()
					if group == nil {
						continue
					}
					groupID := GetGroupUID(group)
					if groupID == child.ID() {
						childItem = group
					} else if groupID == parent.ID() {
						parentItem = group
					}
					// Stop loop if we found everything
					if parentItem != nil && childItem != nil {
						break
					}
				}
				// Create and add link
				link := AddLink(parentItem, childItem)
				scene.AddItem(link.line)
				scene.AddItem(link.dir)
			}
		}
	}

	// Setup drag-and-drop
	view.SetAcceptDrops(true)
	view.SetAlignment(core.Qt__AlignTop | core.Qt__AlignLeft)
	view.ConnectDragMoveEvent(func(event *gui.QDragMoveEvent) {
		if event.Source() != nil && event.Source().IsWidgetType() {
			event.AcceptProposedAction()
		}
	})
	// What item we're currently moving, if any
	var movingItem *widgets.QGraphicsItemGroup
	// Start position of link
	var linkStart *widgets.QGraphicsItemGroup

	itemSize := 64
	view.ConnectDropEvent(func(event *gui.QDropEvent) {
		pos := view.MapToScene(event.Pos())

		// Add item to database
		// For now, we assume all items are requirements
		db := currentProject.Data()
		defer db.Close()
		uid, err := db.AddEmptyRequirement()
		if err != nil {
			widgets.QMessageBox_Warning(
				window, "Failed to add item", err.Error(),
				widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			return
		}
		// Snap to grid
		gridPos := SnapToGrid(pos.ToPoint())
		// Set size and position, TODO: Temporary, not needed later
		req := NewRequirement(uid)
		req.SetPos(gridPos.Y(), gridPos.Y())
		req.SetSize(itemSize*2, itemSize)
		// Add item to view
		scene.AddItem(NewGraphicsItem(
			fmt.Sprintf("%x", uid), gridPos.X(), gridPos.Y(), itemSize*2, itemSize, uid))
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
			if linkBtn.IsChecked() {
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
			// Update item position
			movingItem.SetPos(view.MapToScene(SnapToGrid(event.Pos())))
			// Update link
			UpdateLinkPos(movingItem, movingItem.Pos().X(), movingItem.Pos().Y())
		}
	})
	view.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		if event.Button() == core.Qt__RightButton && view.ItemAt(event.Pos()).Group() != nil {
			// When right clicking item, show edit/delete options
			pos := event.Pos()
			menu := widgets.NewQMenu(nil)
			// Edit option
			menu.AddAction2(gui.QIcon_FromTheme("document-edit"), "Edit").
				ConnectTriggered(func(checked bool) {
					if editWidget, ok := CreateEditWidgetFromPos(pos); ok {
						window.AddDockWidget(core.Qt__RightDockWidgetArea, editWidget)
					}
				})
			// Delete option
			menu.AddAction2(gui.QIcon_FromTheme("delete"), "Delete").
				ConnectTriggered(func(checked bool) {
					// Hakke write ur delete here
				})
			// Show menu at cursor
			menu.Popup(view.MapToGlobal(event.Pos()), nil)
			return
		}

		// We released a button while moving an item
		if movingItem != nil {
			pos := movingItem.Pos()
			// Update link if needed
			// Error handling is already taken care of in UpdateLinkPos
			UpdateLinkPos(movingItem, pos.X(), pos.Y())
			// Update position in database
			// TODO: Assuming Requirement
			NewItem(GetGroupUID(movingItem), TypeRequirement).SetPos(int(pos.X()), int(pos.Y()))
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
			// Create and add link
			link := AddLink(linkStart, view.ItemAt(event.Pos()).Group())
			scene.AddItem(link.line)
			scene.AddItem(link.dir)
			linkStart = nil
		}
	})
	return view
}

func GetGroupUID(group *widgets.QGraphicsItemGroup) int64 {
	return group.Data(0).ToLongLong(nil)
}

func AddLink(parent, child *widgets.QGraphicsItemGroup) Line {
	// Check if map needs to be created
	if links == nil {
		links = make(map[int64][]*Line)
	}
	// Add to database
	db := currentProject.Data()
	defer db.Close()
	// TODO: Assuming requirement
	if err := db.AddItemChild(
		NewItem(GetGroupUID(parent), TypeRequirement), NewItem(GetGroupUID(child), TypeRequirement)); err != nil {
			fmt.Println("error: failed to add link to database:", err)
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
		CreateTriangle(line.Line().Center(), line.Line().Angle()),
	}
	links[parentID] = append(links[parentID], &lineData)
	links[childID] = append(links[childID], &lineData)
	// Return the graphics line to add to scene
	return lineData
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
		// Update direction
		center := l.line.Line().Center()
		l.dir.SetPos2(center.X()-8, center.Y()-8)
		l.dir.SetRotation((-l.line.Line().Angle()) - 90)
	}
}

func NewGraphicsItem(text string, x, y, width, height int, uid int64) *widgets.QGraphicsItemGroup {
	group := widgets.NewQGraphicsItemGroup(nil)
	textItem := widgets.NewQGraphicsTextItem2(text, nil)
	textItem.SetZValue(15)
	shapeItem := widgets.NewQGraphicsRectItem3(0, 0, float64(width), float64(height), nil)
	shapeItem.SetBrush(gui.NewQBrush3(backgroundColor, 1))
	group.AddToGroup(textItem)
	group.AddToGroup(shapeItem)
	group.SetPos2(float64(x), float64(y))
	group.SetData(0, core.NewQVariant1(uid))
	group.SetZValue(10)
	return group
}

// CreateTriangle creates a new 16x16 triangle pointing downwards
func CreateTriangle(pos *core.QPointF, angle float64) *widgets.QGraphicsPolygonItem {
	// Total width/height for triangle
	const size = 16
	// Create each point
	points := []*core.QPointF{
		core.NewQPointF3(0, 0),
		core.NewQPointF3(size, 0),
		core.NewQPointF3(size>>1, size),
	}
	// Create polygon and return it
	poly := widgets.NewQGraphicsPolygonItem2(gui.NewQPolygonF3(points), nil)
	poly.SetPos2(pos.X()-(size>>1), pos.Y()-(size>>1))
	poly.SetPen(gui.NewQPen3(gui.NewQColor3(0, 255, 0, 255)))
	poly.SetTransformOriginPoint2(size>>1, size>>1)
	poly.SetRotation((-angle) - 90)
	return poly
}
