package main

import (
	"crypto/md5"
	"errors"
	"fmt"
)

// Requirement with ItemProprties
type Requirement struct {
	Item
	id int
}

// SaveChanges saves all changes to database
func (req *Requirement) SaveChanges() error {
	return errors.New("error: not implemented")
}

// GetChildren get all children of requirement
func (req Requirement) GetChildren() []Item {
	return nil
}

// RemoveChild remove specified child of requirement
func (req Requirement) RemoveChild(child Item) {
}

// GetHash get hash of requirement
func (req *Requirement) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v", req)))
	fmt.Println("req hash of '", req, "': ", h)
	return h
}

// GetValue gets a value from the database
func (req *Requirement) GetValue(name string) interface{} {
	db := currentProject.GetData()
	defer db.Close()
	return db.GetItemValue(req.GetId(), "Requirements", name)
}

// SetValue sets a value to the database
func (req *Requirement) SetValue(name string, value interface{}) {
	db := currentProject.GetData()
	defer db.Close()
	db.SetItemValue(req.GetId(), "Requirements", name, value)
}

// GetRationale gets the rationale property of the requirement
func (req *Requirement) GetRationale() string {
	return req.GetValue("rationale").(string)
}

// GetFitCriterion of Requirement
func (req *Requirement) GetFitCriterion() string {
	return req.GetValue("fitCriterion").(string)
}

// GetId gets the row ID in the database
func (req Requirement) GetId() int {
	return req.id
}

// GetUid gets the row Uid in the database
func (req Requirement) GetUid() int64 {
	return req.GetValue("uid").(int64)
}

// SetUid sets the Uid in the database
func (req Requirement) SetUid(uid int64) {
	req.SetValue("uid", uid)
}

// GetVersion of Requirement
func (req Requirement) GetVersion() int {
	return req.GetValue("version").(int)
}

// GetShown gets the root as hidden or shown
func (req Requirement) GetShown() bool {
	return req.GetValue("shown").(bool)
}

// SetShown sets the root as hidden or shown
func (req Requirement) SetShown(shown bool) {
	req.SetValue("shown", shown)
}

// GetDescription gets the description from the database
func (req Requirement) GetDescription() string {
	return req.GetValue("description").(string)
}

// AddChild
func (req Requirement) AddChild(child Item) {

}
