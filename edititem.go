package main

import (
	"fmt"

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
		"format-text-bold",
		"format-text-italic",
		"format-text-underline",
		"format-text-strikethrough",
	}

	for _, button := range buttons {
		toolBar.AddAction2(gui.QIcon_FromTheme(button), button).SetCheckable(true)
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
	// Item UID (for debugging)
	layout.AddWidget(CreateGroupBox("Item UID",
		widgets.NewQLabel2(fmt.Sprintf("%v", item.ID()), nil, 0)), 0, 0)

	// Requirement/solution selection
	itemType := widgets.NewQGroupBox2("Item Type", nil)
	reqRadio := widgets.NewQRadioButton2("Problem", nil)
	reqRadio.SetChecked(GetItemType(item) == TypeRequirement)
	solRadio := widgets.NewQRadioButton2("Solution", nil)
	solRadio.SetChecked(GetItemType(item) == TypeSolution)
	itemTypeLayout := widgets.NewQHBoxLayout()
	itemTypeLayout.AddWidget(reqRadio, 1, 0)
	itemTypeLayout.AddWidget(solRadio, 1, 0)
	itemType.SetLayout(itemTypeLayout)
	layout.AddWidget(itemType, 0, 0)

	textOptions := [3]*widgets.QToolBar{}
	textEdits := [3]*widgets.QTextEdit{}
	// TODO: Assume requirement for now
	req, _ := item.(Requirement)
	// Get default values
	textValues := []string{
		req.Description(),
		req.Rationale(),
		req.FitCriterion(),
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
		layout.AddWidget(CreateGroupBox(titles[i], textOptions[i], textEdits[i]), 1, 0)
	}

	// Dock for button connections
	dock := widgets.NewQDockWidget("Edit Item", nil, 0)

	// Button container
	buttons := widgets.NewQHBoxLayout()
	// Save button
	save := widgets.NewQPushButton2("Save", nil)
	save.ConnectReleased(func() {
		// Save description to database and recreate group
		// TODO: Probably not the best solution, but it works
		item.SetDescription(textEdits[Description].ToHtml())
		scene.AddItem(NewGraphicsItem(
			fmt.Sprintf("%v%v", item.ID(), textEdits[Description].ToHtml()),
			int(group.X()), int(group.Y()), 128, 64, item.ID()))
		scene.RemoveItem(group)
		// Save other properties
		req.SetRationale(textEdits[Rationale].ToHtml())
		req.SetFitCriterion(textEdits[FitCriterion].ToHtml())
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
