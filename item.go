package main

type Item struct {
	id int
	version int
	shown bool
	uid int64
	description string
}

func (item *Item) SaveChanges() error {
	return nil
}

func (item *Item) GetHash() [16]byte {
	return [16]byte{}
}