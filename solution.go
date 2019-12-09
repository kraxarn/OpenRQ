package main

import (
	"fmt"
	"crypto/md5"
	"errors"
)

type Link struct {
	name string
	color uint
}

func (link *Link) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v%v", link.name, link.color)))
	fmt.Printf("%x", h)
	return h
}

type Solution struct {
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