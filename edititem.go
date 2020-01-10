package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// TextFormat enum (bold, italic, underline, strikethrough)
type TextFormat int8
const (
	FormatBold          TextFormat = 0
	FormatItalic        TextFormat = 1
	FormatUnderline     TextFormat = 2
	FormatStrikeThrough TextFormat = 3
)

// EntryType enum (description, rationale, fit criterion)
type EntryType int8
const (
	Description  EntryType = 0
	Rationale    EntryType = 1
	FitCriterion EntryType = 2
)

// CreateGroupBox creates a new group box with children in a vertical layout
func CreateGroupBox(title string, children ...widgets.QWidget_ITF) *widgets.QGroupBox {
	layout := widgets.NewQVBoxLayout()
	layout.SetSpacing(0)
	layout.SetContentsMargins(0, 0, 0, 0)
	for i, child := range children {
		layout.AddWidget(child, i, 0)
	}
	groupBox := widgets.NewQGroupBox2(title, nil)
	groupBox.SetLayout(layout)
	return groupBox
}

// CreateTextOptions creates the buttons for the various formatting options
func CreateTextOptions() *widgets.QToolBar {
	toolBar := widgets.NewQToolBar2(nil)

	buttons := []string{
		"format-bold",
		"format-italic",
		"format-underline",
		"format-strikethrough",
	}

	for _, button := range buttons {
		toolBar.AddAction2(GetIcon(button), button).SetCheckable(true)
	}

	return toolBar
}

// MergeFormat formats text in a text view
func MergeFormat(textEdit *widgets.QTextEdit, format *gui.QTextCharFormat) {
	cursor := textEdit.TextCursor()
	if !cursor.HasSelection() {
		cursor.Select(gui.QTextCursor__WordUnderCursor)
	}
	cursor.MergeCharFormat(format)
	textEdit.MergeCurrentCharFormat(format)
}

// CreateEditWidget creates the main window for editing an item
func CreateEditWidget(item Item, group *widgets.QGraphicsItemGroup, scene *widgets.QGraphicsScene) *widgets.QDockWidget {
	// Main vertical layout
	layout := widgets.NewQVBoxLayout()
	// Requirement/solution selection
	itemType := GetItemType(item)
	itemTypeGroup := widgets.NewQGroupBox2("Item Type", nil)
	reqRadio := widgets.NewQRadioButton2("Problem", nil)
	reqRadio.SetChecked(itemType == TypeRequirement)
	solRadio := widgets.NewQRadioButton2("Solution", nil)
	solRadio.SetChecked(itemType == TypeSolution)
	itemTypeLayout := widgets.NewQHBoxLayout()
	itemTypeLayout.AddWidget(reqRadio, 1, 0)
	itemTypeLayout.AddWidget(solRadio, 1, 0)
	itemTypeGroup.SetLayout(itemTypeLayout)
	layout.AddWidget(itemTypeGroup, 0, 0)
	// Label for warning about item type
	itemTypeWarn := widgets.NewQLabel2("Changing item type may cause data loss", nil, 0)
	palette := itemTypeWarn.Palette()
	palette.SetColor2(itemTypeWarn.ForegroundRole(), gui.NewQColor2(core.Qt__red))
	itemTypeWarn.SetPalette(palette)
	itemTypeWarn.Hide()
	layout.AddWidget(itemTypeWarn, 0, 0)
	// Hide/show warning when changing type
	reqRadio.ConnectReleased(func() {
		itemTypeWarn.SetVisible(itemType == TypeSolution)
	})
	solRadio.ConnectReleased(func() {
		itemTypeWarn.SetVisible(itemType == TypeRequirement)
	})

	textOptions := [3]*widgets.QToolBar{}
	textEdits := [3]*widgets.QTextEdit{}
	textGroups := [4]*widgets.QGroupBox{}

	// Hide/show when clicking requirement/solution
	reqRadio.ConnectReleased(func() {
		textGroups[1].Show()
		textGroups[2].Show()
		textGroups[3].Hide()
	})
	solRadio.ConnectReleased(func() {
		textGroups[1].Hide()
		textGroups[2].Hide()
		textGroups[3].Show()
	})

	// Get default values
	req, isReq := item.(Requirement)
	textValues := [3]string{
		item.Description(),
	}
	// Also set rationale and fit criterion if requirement
	if isReq {
		textValues[1] = req.Rationale()
		textValues[2] = req.FitCriterion()
	}

	for i := 0; i < len(textOptions); i++ {
		t := CreateTextOptions()
		t.SetVisible(false)
		textOptions[i] = t
		// Attack tool bar buttons for actions
		i2 := i
		for format, action := range t.Actions() {
			f := format
			action.ConnectTriggered(func(checked bool) {
				switch TextFormat(f) {
				// Bold text
				case FormatBold:
					charFormat := gui.NewQTextCharFormat()
					fontWeight := gui.QFont__Normal
					if checked {
						fontWeight = gui.QFont__Bold
					}
					charFormat.SetFontWeight(int(fontWeight))
					MergeFormat(textEdits[i2], charFormat)
				// Italic text
				case FormatItalic:
					charFormat := gui.NewQTextCharFormat()
					charFormat.SetFontItalic(checked)
					MergeFormat(textEdits[i2], charFormat)
				// Strike through text
				case FormatStrikeThrough:
					charFormat := gui.NewQTextCharFormat()
					charFormat.SetFontStrikeOut(checked)
					MergeFormat(textEdits[i2], charFormat)
				// Underlined text
				case FormatUnderline:
					charFormat := gui.NewQTextCharFormat()
					charFormat.SetFontUnderline(checked)
					MergeFormat(textEdits[i2], charFormat)
				}
			})
		}
	}

	updateTextOptions := func(index int) {
		for _, option := range textOptions {
			option.SetVisible(false)
		}
		textOptions[index].SetVisible(true)
	}

	// Looping through Description, Rationale, Fit Criterion.
	titles := []string{
		"Description",
		"Rationale",
		"Fit Criterion",
	}
	for i := 0; i < len(titles); i++ {
		textEdits[i] = widgets.NewQTextEdit(nil)
		textEdits[i].SetHtml(textValues[i])
		// Local copy of i
		i2 := i
		// Show/hide font options on selection
		textEdits[i].ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
			updateTextOptions(i2)
		})
		// Update font options when selecting new text
		textEdits[i].ConnectCurrentCharFormatChanged(func(charFormat *gui.QTextCharFormat) {
			actions := textOptions[i2].Actions()
			font := charFormat.Font()
			actions[FormatBold].SetChecked(font.Bold())
			actions[FormatItalic].SetChecked(font.Italic())
			actions[FormatUnderline].SetChecked(font.Underline())
			actions[FormatStrikeThrough].SetChecked(font.StrikeOut())
		})
		// Add text edit in group box to main layout
		textGroups[i] = CreateGroupBox(titles[i], textOptions[i], textEdits[i])
		layout.AddWidget(textGroups[i], 1, 0)
	}

	// Add image selection for solution
	// TODO: Temporary placeholder text
	noImagesText := widgets.NewQLabel2("No attachments", nil, core.Qt__Widget)
	noImagesText.SetEnabled(false)
	textGroups[3] = CreateGroupBox("Attachments", noImagesText)
	layout.AddWidget(textGroups[3], 0, 0)

	// Hide stuff
	if itemType == TypeRequirement {
		textGroups[3].Hide()
	} else {
		textGroups[1].Hide()
		textGroups[2].Hide()
	}

	// Dock for button connections
	dock := widgets.NewQDockWidget(fmt.Sprintf("Edit Item (%v)", item.ToString()), nil, 0)
	// Button container
	buttons := widgets.NewQHBoxLayout()
	// Save button
	save := widgets.NewQPushButton2("Save", nil)
	// Temporary values
	var itemX, itemY, itemW, itemH int
	save.ConnectReleased(func() {
		// Check if we are changing item type
		changingType := itemTypeWarn.IsVisible()
		if changingType {
			// We are changing type, we need to delete and add item again
			db := currentProject.Data()
			defer db.Close()
			// Save old properties
			oldItem := item
			itemUID := item.UID()
			itemX, itemY = item.Pos()
			itemW, itemH = item.Size()
			itemParent := item.Parent()
			// Get links and delete the old ones
			itemLinks := links[item]
			delete(links, item)
			// Delete old item
			if err := db.RemoveItem(item); err != nil {
				fmt.Println("error: failed to delete old item:", err)
				return
			}
			// Add new as set item type
			if reqRadio.IsChecked() {
				// Add new requirement with set values
				newID, err := db.AddEmptyRequirement()
				if err != nil {
					fmt.Println("error: failed to create new requirement:", err)
					return
				}
				// Create new and set UID, rest is same as updating item
				item = NewRequirement(newID)
			} else {
				// Add new requirement with set values
				newID, err := db.AddEmptySolution()
				if err != nil {
					fmt.Println("error: failed to create new solution:", err)
					return
				}
				// Create new and set UID, rest is same as updating item
				item = NewSolution(newID)
			}
			// Set UID and update and update requirement cast
			// (rest is the same as updating an item)
			item.SetUID(itemUID)
			// Also set position and size same as old
			item.SetPos(itemX, itemY)
			item.SetSize(itemW, itemH)
			// Update any links
			if err := db.UpdateItemChildren(oldItem, item); err != nil {
				fmt.Println("warning: failed to update item children:", err)
			}
			// Update parent
			if itemParent != nil {
				item.SetParent(itemParent)
			}
			// Update links
			for _, link := range itemLinks {
				if link.parent == oldItem {
					link.parent = item
				} else if link.child == oldItem {
					link.child = item
				}
			}
			links[item] = itemLinks
			// Update requirement values
			req, isReq = item.(Requirement)
		}

		// Properties both items need
		item.SetDescription(textEdits[Description].ToHtml())
		// Requirements also need rationale and fit criterion updated
		if isReq {
			req.SetRationale(textEdits[Rationale].ToHtml())
			req.SetFitCriterion(textEdits[FitCriterion].ToHtml())
		}
		// Recreate group with new item
		scene.AddItem(NewGraphicsItem(textEdits[Description].ToHtml(),
			int(group.X()), int(group.Y()), 128, 64, item))
		scene.RemoveItem(group)
		// Close window
		dock.Close()
	})
	buttons.AddWidget(save, 1, 0)
	// Discard button
	discard := widgets.NewQPushButton2("Discard", nil)
	discard.ConnectReleased(func() {
		dock.Close()
	})
	buttons.AddWidget(discard, 1, 0)
	layout.AddLayout(buttons, 0)

	// Put layout in a widget
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)

	// Set dock to the created widget and return it
	dock.SetWidget(widget)
	return dock
}
