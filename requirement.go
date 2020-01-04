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

func (req Requirement) ID() int64 {
	return req.id
}

func (req Requirement) UID() int64 {
	return req.GetValueInt64("uid")
}

func (req Requirement) SetUID(uid int64) {
	req.SetValue("uid", uid)
}

func (req Requirement) Description() string {
	return req.GetValueString("description")
}

func (req Requirement) SetDescription(value string)  {
	req.SetValue("description", value)
}

func (req *Requirement) Rationale() string {
	return req.GetValueString("rationale")
}

func (req *Requirement) SetRationale(value string)  {
	req.SetValue("rationale", value)
}

func (req *Requirement) FitCriterion() string {
	return req.GetValueString("fitCriterion")
}

func (req *Requirement) SetFitCriterion(value string)  {
	req.SetValue("fitCriterion", value)
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

func (req Requirement) Parent() Item {
	// Check if item has parent
	if req.IsPropertyNull("parent") {
		return nil
	}
	// If parent is not null, assume parent type isn't null either
	return NewItem(req.GetValueInt64("parent"), ItemType(req.GetValueInt("parentType")))
}

func (req Requirement) SetParent(parent Item) {
	if parent == nil {
		req.SetValues(map[string]interface{}{
			"parent":     nil,
			"parentType": nil,
		})
	} else {
		req.SetValues(map[string]interface{}{
			"parent":     parent.ID(),
			"parentType": GetItemType(parent),
		})
	}
}

func (req Requirement) Hash() [16]byte {
	return md5.Sum([]byte(fmt.Sprintf("%v", req)))
}

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

func (req Requirement) IsPropertyNull(columnName string) bool {
	db := currentProject.Data()
	defer db.Close()
	return db.IsItemPropertyNull(req.id, "Requirements", columnName)
}