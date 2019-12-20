package main

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type TextFormat int8

const (
	FormatBold          TextFormat = 0
	FormatItalic        TextFormat = 1
	FormatUnderline     TextFormat = 2
	FormatStrikeThrough TextFormat = 3
)

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

func CreateEditWidget() *widgets.QDockWidget {
	// Main vertical layout
	layout := widgets.NewQVBoxLayout()

	// Requirement/solution selection
	itemType := widgets.NewQGroupBox2("Item Type", nil)
	reqRadio := widgets.NewQRadioButton2("Requirement", nil)
	solRadio := widgets.NewQRadioButton2("Solution", nil)
	itemTypeLayout := widgets.NewQHBoxLayout()
	itemTypeLayout.AddWidget(reqRadio, 1, 0)
	itemTypeLayout.AddWidget(solRadio, 1, 0)
	itemType.SetLayout(itemTypeLayout)
	layout.AddWidget(itemType, 0, 0)

	textOptions := [3]*widgets.QToolBar{}

	for i := 0; i < 3; i++ {
		t := CreateTextOptions()
		t.SetVisible(false)
		textOptions[i] = t
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
		textEdit := widgets.NewQTextEdit(nil)
		// Local copy of i
		i2 := i
		textEdit.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
			updateTextOptions(i2)
		})
		layout.AddWidget(CreateGroupBox(titles[i], textOptions[i], textEdit), 1, 0)
	}

	// Save and dismiss
	layout.AddWidget(widgets.NewQPushButton2("Save", nil), 1, 0)

	// Put layout in a widget
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	// Set dock to the created widget and return it
	dock := widgets.NewQDockWidget("Edit Item", nil, 0)
	dock.SetWidget(widget)
	return dock
}
