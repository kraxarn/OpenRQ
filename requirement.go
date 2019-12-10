package main

import (
	"crypto/md5"
	"errors"
	"fmt"
)

// Requirement with ItemProprties
type Requirement struct {
	Item
	ItemProperties
	rationale    string
	fitCriterion string
}

// SaveChanges saves all changes to database
func (req *Requirement) SaveChanges() error {
	return errors.New("error: not implemented")
}

// GetChildren get all children of requirement
func (req *Requirement) GetChildren() []Item {
	return nil
}

// RemoveChild remove specified child of requirement
func (req *Requirement) RemoveChild() []Item {
	return nil
}

// GetHash get hash of requirement
func (req *Requirement) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v", req)))
	fmt.Println("req hash of '", req, "': ", h)
	return h
}

// Rationale of Requirement
func (req *Requirement) GetRationale() string {
	return req.rationale
}

// FitCriterion of Requirement
func (req *Requirement) GetFitCriterion() string {
	return req.fitCriterion
}

func (req Requirement) GetId() int {
	return req.id
}

func (req Requirement) GetUid() int64 {
	return req.uid
}

func (req Requirement) SetUid(uid int64) {
	req.uid = uid
}

func (req Requirement) GetVersion() int {
	return req.version
}

func (req Requirement) GetShown() bool {
	return req.shown
}

func (req Requirement) SetShown(shown bool) {
	req.shown = shown
}

func (req Requirement) GetDescription() string {
	return req.description
}
