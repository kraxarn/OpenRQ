package main

// Item that is either a solutions or a requirement
type Item interface {
	GetId() int64

	GetUid() int64
	SetUid(uid int64)

	GetVersion() int

	GetShown() bool
	SetShown(shown bool)

	GetDescription() string
	SetDescription(description string)

	Pos() (int, int)
	SetPos(x, y int)

	Size() (int, int)
	SetSize(w, h int)

	AddChild(child Item)
	RemoveChild(child Item)

	GetChildren() []Item
	Parent() (parentID int64, parentType ItemType, found bool)
}

func NewItem(id int64, itemType ItemType) Item {
	var item Item
	if itemType == TypeRequirement {
		item = NewRequirement(id)
	} else {
		item = NewSolution(id)
	}
	return item
}