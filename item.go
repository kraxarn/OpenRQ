package main

import (
	"encoding/json"
	"fmt"
)

// Item that is either a solutions or a requirement
type Item interface {
	json.Marshaler

	ID() int64

	UID() int64
	SetUID(uid int64)

	Version() int

	Hash() [16]byte

	Shown() bool
	SetShown(shown bool)

	Description() string
	SetDescription(description string)

	Pos() (int, int)
	SetPos(x, y int)

	Size() (int, int)
	SetSize(w, h int)

	AddChild(child Item)
	RemoveChild(child Item)
	Children() []Item

	Parent() Item
	SetParent(parent Item)

	IsPropertyNull(columnName string) bool
	ToString() string
}

func NewItem(id int64, itemType ItemType) Item {
	var item Item
	if itemType == TypeRequirement {
		item = NewRequirement(id)
	} else if itemType == TypeSolution {
		item = NewSolution(id)
	} else {
		fmt.Println("error: failed to create item from id", id, "type", itemType)
	}
	return item
}