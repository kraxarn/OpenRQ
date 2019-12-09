package main

import (
	"fmt"
	"crypto/md5"
	"errors"
)

type Requirement struct {
	Item
	rationale string
	fitCriterion string
}

func (req *Requirement) SaveChanges() error {
	return errors.New("error: not implemented")
}

func (req *Requirement)	GetChildren() []Item {
	return nil
}

func (req *Requirement) RemoveChild() []Item {
	return nil
}

func (req *Requirement) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v", req)))
	fmt.Println("req hash of '", req, "': ", h)
	return h
}