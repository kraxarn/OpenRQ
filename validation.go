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
	// Final returned items
	items = make([]Item, 0)
	// Items we have already added
	added := map[Item]int{}
	// Loop through all items in view
	for _, item := range view.Items() {
		// Get group and make sure it's valid
		group := item.Group()
		if group == nil || group.Type() == 0 {
			continue
		}
		// Get item
		groupItem := GetGroupItem(group)
		// Ignore if already added
		if ContainsItem(added, groupItem) {
			continue
		}
		added[groupItem] = 0
		// Get children count and if it has a parent (not a root)
		children := 0
		hasParent := false
		for _, link := range links[groupItem] {
			if link.parent == groupItem {
				children++
			} else if link.child == groupItem {
				hasParent = true
				break
			}
		}
		// If it had a parent, it's not a root
		if hasParent {
			continue
		}
		// Otherwise, add if it had more than one child
		if children > 1 {
			items = append(items, groupItem)
		}
	}
	return items
}

// Validates tree to check if there are any links to each other
func ValidateLoops() (items []Item) {
	// Final returned splice
	items = make([]Item, 0)
	// Items we have already added
	added := map[Item]int{}
	// Loop through all links
	for _, link := range links {
		// Loop through all links related to that item
		for _, l1 := range link {
			// Loop through all links again checking for opposite links
			for _, l2 := range link {
				// They are linked to each other
				if l1.parent == l2.child && l1.child == l2.parent && !ContainsItem(added, l1.child) {
					items = append(items, l1.child)
					added[l1.child] = 0
				}
			}
		}
	}
	return items
}

func ValidateLinkErrors() (items []Item) {
	// Final returned items
	items = make([]Item, 0)
	// Items we have already added
	added := map[Item]int{}
	// Loop through all items in view
	for _, item := range view.Items() {
		// Get group and make sure it's valid
		group := item.Group()
		if group == nil || group.Type() == 0 {
			continue
		}
		// Get item
		groupItem := GetGroupItem(group)
		// Ignore if already added
		if ContainsItem(added, groupItem) {
			continue
		}
		// Get children count and if it has a parent (not a root)
		parents := 0
		for _, link := range links[groupItem] {
			if link.child == groupItem {
				parents++
			}
		}
		// If it has more than one parent, something is wrong
		if parents > 1 {
			items = append(items, groupItem)
			added[groupItem] = 0
		}
	}
	return items
}

type ValidationOption int8
const (
	SameType 	ValidationOption = 0
	OneRoot		ValidationOption = 1
	LinkLoop	ValidationOption = 2
	LinkError	ValidationOption = 3
)

type ValidationResult string
const (
	ValidateOK       ValidationResult = "validate-ok"
	ValidateFail     ValidationResult = "validate-fail"
	ValidateDisabled ValidationResult = "validate-disabled"
	ValidateNone     ValidationResult = "validate-none"
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
	case LinkError:
		text = "Invalid links"
		info = "Link that could not be saved to the project file"
	}
	item := widgets.NewQListWidgetItem3(GetIcon(string(result)), text, nil, 0)
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
		true, true, true, true,
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
		// Start validation timer
		start := time.Now()
		// Run link validation
		if enabled[SameType] {
			valLinks := ValidateLinks()
			for _, item := range valLinks {
				items.AddItem(fmt.Sprintf("%v %v\n(links to same type)", GetItemName(item), item.ID()))
			}
			results.Item(int(SameType)).SetIcon(GetIcon(string(GetValidationResult(len(valLinks)))))
		}
		// Run root validation
		if enabled[OneRoot] {
			valRoots := ValidateRoots()
			for _, item := range valRoots {
				items.AddItem(fmt.Sprintf("%v %v\n(one-to-one root)", GetItemName(item), item.ID()))
			}
			results.Item(int(OneRoot)).SetIcon(GetIcon(string(GetValidationResult(len(valRoots)))))
		}
		// Run loop validation
		if enabled[LinkLoop] {
			valLoops := ValidateLoops()
			for _, item := range valLoops {
				items.AddItem(fmt.Sprintf("%v %v\n(linking loop)", GetItemName(item), item.ID()))
			}
			results.Item(int(LinkLoop)).SetIcon(GetIcon(string(GetValidationResult(len(valLoops)))))
		}
		// Run invalid link validation
		if enabled[LinkError] {
			valErrors := ValidateLinkErrors()
			for _, item := range valErrors {
				items.AddItem(fmt.Sprintf("%v %v\n(invalid link)", GetItemName(item), item.ID()))
			}
			results.Item(int(LinkError)).SetIcon(GetIcon(string(GetValidationResult(len(valErrors)))))
		}
		// Enable them again
		runBtn.SetText("Run now")
		runBtn.SetEnabled(true)
		results.SetEnabled(true)
		items.SetEnabled(true)
		// Set last validation time
		title.SetTitle(fmt.Sprintf("Last validation: %v (%v ms)", time.Now().Format("15:04"), time.Now().Sub(start).Milliseconds()))
	})
	layout.AddWidget(runBtn, 0, 0)

	// Create initial results
	for i, e := range enabled {
		results.AddItem2(CreateValidationResult(ValidationOption(i), GetDefaultValidationResult(e)))
	}
	layout.AddWidget(title, 0, 0)
	// Show menu when clicking on result item
	results.ConnectItemPressed(func(item *widgets.QListWidgetItem) {
		i := results.Row(item)
		menu := widgets.NewQMenu(nil)
		action := menu.AddAction("Enabled")
		action.SetCheckable(true)
		action.SetChecked(enabled[i])
		action.ConnectTriggered(func(checked bool) {
			enabled[i] = action.IsChecked()
			item.SetIcon(GetIcon(string(GetDefaultValidationResult(enabled[i]))))
		})
		menu.Popup(gui.QCursor_Pos(), nil)
	})
	// Create list to show affected items
	itemGroup := CreateGroupBox("Affected Items", items)
	layout.AddWidget(itemGroup, 1, 0)

	// Convert layout to widget and return it
	widget := widgets.NewQWidget(nil, core.Qt__Widget)
	widget.SetLayout(layout)
	widget.SetMaximumWidth(250)
	widget.SetMinimumWidth(175)
	return widget
}