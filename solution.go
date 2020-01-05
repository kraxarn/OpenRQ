package main

import (
	"crypto/md5"
	"fmt"
	"os"
)

type Solution struct {
	Item
	id int64
}

func NewSolution(id int64) Solution {
	sol := Solution{}
	sol.id = id
	if id == 0 {
		fmt.Fprintln(os.Stderr, "warning: trying to create solution with id 0")
	}
	return sol
}

func (sol Solution) ID() int64 {
	return sol.id
}

func (sol Solution) UID() int64 {
	return sol.GetValueInt64("uid")
}

func (sol Solution) SetUID(uid int64) {
	sol.SetValue("uid", uid)
}

func (sol Solution) Description() string {
	return sol.GetValueString("description")
}

func (sol Solution) SetDescription(value string) {
	sol.SetValue("description", value)
}

func (sol Solution) Pos() (int, int) {
	var x, y int
	sol.GetValues(map[string]interface{}{
		"x": &x,
		"y": &y,
	})
	return x, y
}

func (sol Solution) SetPos(x, y int) {
	sol.SetValues(map[string]interface{} {
		"x": x,
		"y": y,
	})
}

func (sol Solution) Size() (int, int) {
	var width, height int
	sol.GetValues(map[string]interface{}{
		"width": &width,
		"height": &height,
	})
	return width, height
}

func (sol Solution) SetSize(w, h int) {
	sol.SetValues(map[string]interface{} {
		"width": w,
		"height": h,
	})
}

func (sol Solution) Parent() Item {
	// Check if item has parent
	if sol.IsPropertyNull("parent") {
		return nil
	}
	// If parent is not null, assume parent type isn't null either
	return NewItem(sol.GetValueInt64("parent"), ItemType(sol.GetValueInt("parentType")))
}

func (sol Solution) SetParent(parent Item) {
	if parent == nil {
		sol.SetValues(map[string]interface{}{
			"parent":     nil,
			"parentType": nil,
		})
	} else {
		sol.SetValues(map[string]interface{}{
			"parent":     parent.ID(),
			"parentType": GetItemType(parent),
		})
	}
}

func (sol Solution) Hash() [16]byte {
	return md5.Sum([]byte(fmt.Sprintf("%v", sol)))
}

func (sol *Solution) GetValue(name string, value interface{}) {
	sol.GetValues(map[string]interface{}{
		name: value,
	})
}

func (sol *Solution) GetValues(nameValues map[string]interface{}) {
	db := currentProject.Data()
	defer db.Close()
	for key, value := range nameValues {
		if err := db.GetItemValue(sol.ID(), "Solutions", key, value); err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed to get property", key, ":", err)
		}
	}
}

func (sol *Solution) GetValueString(name string) string {
	var val string
	sol.GetValue(name, &val)
	return val
}

func (sol *Solution) GetValueInt(name string) int {
	var val int
	sol.GetValue(name, &val)
	return val
}

func (sol *Solution) GetValueInt64(name string) int64 {
	var val int64
	sol.GetValue(name, &val)
	return val
}

func (sol *Solution) SetValue(name string, value interface{}) {
	sol.SetValues(map[string]interface{}{
		name: value,
	})
}

func (sol *Solution) SetValues(nameValues map[string]interface{}) {
	db := currentProject.Data()
	defer db.Close()
	for key, value := range nameValues {
		db.SetItemValue(sol.ID(), "Solutions", key, value)
	}
}

func (sol Solution) IsPropertyNull(columnName string) bool {
	db := currentProject.Data()
	defer db.Close()
	return db.IsItemPropertyNull(sol.id, "Solutions", columnName)
}

func (sol Solution) ToString() string {
	return fmt.Sprintf("solution %v", sol.id)
}