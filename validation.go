package main

import (
	"fmt"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func ContainsItem(items map[Item]int, target Item) bool {
	_, ok := items[target]
	return ok
}

func GetItemName(item Item) string {
	switch GetItemType(item) {
	case TypeRequirement:
		return "Problem"
	case TypeSolution:
		return "Solution"
	}
	return ""
}

// Validates link to check that links are not the same type
func ValidateLinks() (items []Item) {
	// Final returned splice
	items = make([]Item, 0)
	// Items we have already added
	added := map[Item]int{}
	// Loop through all links
	for _, link := range links {
		// Loop through all links related to that item
		for _, l := range link {
			// Check if child has same item type and not already added
			if GetItemType(l.parent) == GetItemType(l.child) && !ContainsItem(added, l.child) {
				items = append(items, l.child)
				added[l.child] = 0
			}
		}
	}
	return items
}

// Validates roots to check if they have a one-to-one relation
func ValidateRoots() (items []Item) {
	items = make([]Item, 0)
	for _, item := range view.Items() {
		group := item.Group()
		if group == nil || group.Type() == 0 {
			continue
		}
		// TODO
	}
	return items
}

// Validates tree to check if there are any links to each other
func ValidateLoops() (items []Item) {
	// TODO
	return nil
}

type ValidationOption int8
const (
	SameType ValidationOption = 0
	OneRoot  ValidationOption = 1
	LinkLoop ValidationOption = 2
)

type ValidationResult string
const (
	ValidateOK       ValidationResult = "emblem-checked"
	ValidateFail     ValidationResult = "emblem-error"
	ValidateDisabled ValidationResult = "emblem-pause"
	ValidateNone     ValidationResult = "emblem-question"
)

func CreateValidationResult(option ValidationOption, result ValidationResult) *widgets.QListWidgetItem {
	var text, info string
	switch option {
	case SameType:
		text = "Links to same type"
		info = "Linking items of the same type"
	case OneRoot:
		text = "One-to-one root"
		info = "Roots of the current tree having an one-to-one relation to its child"
	case LinkLoop:
		text = "Linking loop"
		info = "Items that link to each other in a loop"
	}
	item := widgets.NewQListWidgetItem3(gui.QIcon_FromTheme(string(result)), text, nil, 0)
	item.SetToolTip(info)
	return item
}

func GetDefaultValidationResult(enabled bool) ValidationResult {
	if enabled {
		return ValidateNone
	}
	return ValidateDisabled
}

func GetValidationResult(itemCount int) ValidationResult {
	if itemCount > 0 {
		return ValidateFail
	}
	return ValidateOK
}

// CreateValidationEngineLayout creates the validation engine window
func CreateValidationEngineLayout() *widgets.QWidget {
	// Enable all validations by default
	// (this should maybe be loaded/saved from database)
	enabled := []bool{
		true, true, true,
	}
	// Main vertical box
	layout := widgets.NewQVBoxLayout()
	// List of validation results
	results := widgets.NewQListWidget(nil)
	// List of affected items
	items := widgets.NewQListWidget(nil)
	// Main container for validation results
	title := CreateGroupBox("Last validation: never", results)
	// Option to validate
	runBtn := widgets.NewQPushButton2("Run now", nil)
	runBtn.ConnectReleased(func() {
		// Update text and disable all options
		runBtn.SetText("Running...")
		runBtn.SetEnabled(false)
		results.SetEnabled(false)
		items.SetEnabled(false)
		// Empty list of affected items
		items.Clear()
		// Run link validation
		if enabled[SameType] {
			valLinks := ValidateLinks()
			for _, item := range valLinks {
				items.AddItem(fmt.Sprintf("%v %v", GetItemName(item), item.ID()))
			}
			results.Item(int(SameType)).SetIcon(gui.QIcon_FromTheme(string(GetValidationResult(len(valLinks)))))
		}
		// Enable them again
		runBtn.SetText("Run now")
		runBtn.SetEnabled(true)
		results.SetEnabled(true)
		items.SetEnabled(true)
		// Set last validation time
		title.SetTitle(fmt.Sprintf("Last validation: %v", time.Now().Format("15:04")))
	})
	layout.AddWidget(runBtn, 0, 0)

	/*
	 * List icons:
	 * ok:		emblem-checked
	 * fail:	emblem-error
	 * off:		emblem-pause
	 * none:	emblem-question
	 */

	for i, e := range enabled {
		results.AddItem2(CreateValidationResult(ValidationOption(i), GetDefaultValidationResult(e)))
	}
	layout.AddWidget(title, 1, 0)

	results.ConnectItemPressed(func(item *widgets.QListWidgetItem) {
		i := results.Row(item)
		menu := widgets.NewQMenu(nil)
		action := menu.AddAction("Enabled")
		action.SetCheckable(true)
		action.SetChecked(enabled[i])
		action.ConnectTriggered(func(checked bool) {
			enabled[i] = action.IsChecked()
			item.SetIcon(gui.QIcon_FromTheme(string(GetDefaultValidationResult(enabled[i]))))
		})
		menu.Popup(gui.QCursor_Pos(), nil)
	})

	noItems := widgets.NewQListWidgetItem2("No items to show", nil, 0)
	noItems.SetFlags(0)
	itemGroup := CreateGroupBox("Affected Items", items)
	layout.AddWidget(itemGroup, 1, 0)

	// Convert layout to widget and return it
	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	widget.SetMaximumWidth(250)
	widget.SetMinimumWidth(175)
	return widget
}