package main

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type Link struct {
	name  string
	color uint
}

func (link *Link) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v%v", link.name, link.color)))
	fmt.Printf("%x", h)
	return h
}

type Solution struct {
	ItemProperties
	Item
}

func (sol *Solution) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v", sol)))
	fmt.Println("sol hash of '", sol, "': ", h)
	return h
}

func (sol *Solution) SaveChanges() error {
	return errors.New("error: not implemented")
}

func (sol *Solution) GetChildren() []Item {
	return nil
}

func (sol *Solution) RemoveChild() []Item {
	return nil
}

func (sol Solution) Id() int {
	return sol.id
}

func (sol Solution) Uid() int64 {
	return sol.uid
}

func (sol Solution) Version() int {
	return sol.version
}

func (sol Solution) Shown() bool {
	return sol.shown
}

func (sol Solution) Description() string {
	return sol.description
}
