package main

// Item that is either a solutions or a requirement
type Item interface {
	GetId() int

	GetUid() int64
	SetUid(uid int64)

	GetVersion() int

	GetShown() bool
	SetShown(shown bool)

	GetDescription() string
	SetDescription(description string)
}

// ItemProperties properties (but not methods) for items
type ItemProperties struct {
	id          int
	uid         int64
	version     int
	shown       bool
	description string
}
