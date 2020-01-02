package main

import (
	"crypto/md5"
	"fmt"
	"os"
)

type Requirement struct {
	Item
	id int64
}

func NewRequirement(id int64) Requirement {
	req := Requirement{}
	req.id = id
	if id == 0 {
		fmt.Fprintln(os.Stderr, "warning: trying to create requirement with id 0")
	}
	return req
}

// GetChildren get all children of requirement
func (req Requirement) Children() []Item {
	return nil
}

// RemoveChild remove specified child of requirement
func (req Requirement) RemoveChild(child Item) {
}

// GetHash get hash of requirement
func (req Requirement) Hash() [16]byte {
	return md5.Sum([]byte(fmt.Sprintf("%v", req)))
}

// GetValue gets a value from the database
func (req *Requirement) GetValue(name string, value interface{}) {
	req.GetValues(map[string]interface{}{
		name: value,
	})
}

func (req *Requirement) GetValues(nameValues map[string]interface{}) {
	db := currentProject.Data()
	defer db.Close()
	for key, value := range nameValues {
		if err := db.GetItemValue(req.ID(), "Requirements", key, value); err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed to get property", key, ":", err)
		}
	}
}

func (req *Requirement) GetValueString(name string) string {
	var val string
	req.GetValue(name, &val)
	return val
}

func (req *Requirement) GetValueInt(name string) int {
	var val int
	req.GetValue(name, &val)
	return val
}

func (req *Requirement) GetValueInt64(name string) int64 {
	var val int64
	req.GetValue(name, &val)
	return val
}

// SetValue sets a value to the database
func (req *Requirement) SetValue(name string, value interface{}) {
	req.SetValues(map[string]interface{}{
		name: value,
	})
}

func (req *Requirement) SetValues(nameValues map[string]interface{}) {
	db := currentProject.Data()
	defer db.Close()
	for key, value := range nameValues {
		db.SetItemValue(req.ID(), "Requirements", key, value)
	}
}

// GetRationale gets the rationale property of the requirement
func (req *Requirement) Rationale() string {
	return req.GetValueString("rationale")
}

func (req *Requirement) SetRationale(value string)  {
	req.SetValue("rationale", value)
}

// GetFitCriterion of Requirement
func (req *Requirement) FitCriterion() string {
	return req.GetValueString("fitCriterion")
}

func (req *Requirement) SetFitCriterion(value string)  {
	req.SetValue("fitCriterion", value)
}

// GetId gets the row ID in the database
func (req Requirement) ID() int64 {
	return req.id
}

// GetUid gets the row Uid in the database
func (req Requirement) UID() int64 {
	return req.GetValueInt64("uid")
}

// SetUid sets the Uid in the database
func (req Requirement) SetUID(uid int64) {
	req.SetValue("uid", uid)
}

// GetVersion of Requirement
func (req Requirement) Version() int {
	return req.GetValueInt("version")
}

// GetShown gets the root as hidden or shown
func (req Requirement) Shown() bool {
	var val bool
	req.GetValue("shown", &val)
	return val
}

// SetShown sets the root as hidden or shown
func (req Requirement) SetShown(shown bool) {
	req.SetValue("shown", shown)
}

// GetDescription gets the description from the database
func (req Requirement) Description() string {
	return req.GetValueString("description")
}

func (req Requirement) SetDescription(value string)  {
	req.SetValue("description", value)
}

// AddChild
func (req Requirement) AddChild(child Item) {

}

func (req Requirement) Pos() (int, int) {
	var x, y int
	req.GetValues(map[string]interface{}{
		"x": &x,
		"y": &y,
	})
	return x, y
}

func (req Requirement) SetPos(x, y int) {
	req.SetValues(map[string]interface{} {
		"x": x,
		"y": y,
	})
}

func (req Requirement) Size() (int, int) {
	var width, height int
	req.GetValues(map[string]interface{}{
		"width": &width,
		"height": &height,
	})
	return width, height
}

func (req Requirement) SetSize(w, h int) {
	req.SetValues(map[string]interface{} {
		"width": w,
		"height": h,
	})
}

func (req Requirement) IsPropertyNull(columnName string) bool {
	db := currentProject.Data()
	defer db.Close()
	return db.IsItemPropertyNull(req.id, "Requirements", columnName)
}