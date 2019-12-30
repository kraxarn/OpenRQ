package main

// Item that is either a solutions or a requirement
type Item interface {
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