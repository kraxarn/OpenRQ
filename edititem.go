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
	for _, child := range children {
		layout.AddWidget(child, 1, 0)
	}
	groupBox := widgets.NewQGroupBox2(title, nil)
	groupBox.SetLayout(layout)
	return groupBox
}

func CreateIconButton(icon string) *widgets.QPushButton {
	button := widgets.NewQPushButton(nil)
	button.SetCheckable(true)
	button.SetIcon(gui.QIcon_FromTheme(icon))
	button.SetFlat(true)
	return button
}

func CreateTextOptions() (*widgets.QWidget, []*widgets.QPushButton) {
	layout := widgets.NewQHBoxLayout()

	buttons := []*widgets.QPushButton{
		CreateIconButton("format-text-bold"),
		CreateIconButton("format-text-italic"),
		CreateIconButton("format-text-underline"),
		CreateIconButton("format-text-strikethrough"),
	}

	for _, button := range buttons {
		layout.AddWidget(button, 1, 0)
	}

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	return widget, buttons
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

	textOptions := [3]*widgets.QWidget{}
	textButtons := [3][]*widgets.QPushButton{}

	for i := 0; i < 3; i++ {
		o, b := CreateTextOptions()
		//o.SetVisible(false)
		textOptions[i] = o
		textButtons[i] = b
	}

	updateTextOptions := func(index int) {
		for _, option := range textOptions {
			option.SetVisible(false)
		}
		textOptions[index].SetVisible(true)
	}

	// Description
	descText := widgets.NewQTextEdit(nil)
	descText.ConnectFocusInEvent(func(event *gui.QFocusEvent) {
		updateTextOptions(0)
	})
	layout.AddWidget(CreateGroupBox("Description", textOptions[0], descText), 1, 0)

	// Rationale
	ratText := widgets.NewQTextEdit(nil)
	ratText.ConnectFocusInEvent(func(event *gui.QFocusEvent) {
		updateTextOptions(1)
	})
	layout.AddWidget(CreateGroupBox("Rationale", textOptions[1], ratText), 1, 0)

	// Fit criterion
	fitText := widgets.NewQTextEdit(nil)
	fitText.ConnectFocusInEvent(func(event *gui.QFocusEvent) {
		updateTextOptions(2)
	})
	layout.AddWidget(CreateGroupBox("Fit Criterion", textOptions[2], fitText), 1, 0)

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
