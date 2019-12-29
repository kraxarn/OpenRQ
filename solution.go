package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
)

type Link struct {
	name  string
	color uint
}

// GetHash gets the md5 hash of the item
func (link *Link) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v%v", link.name, link.color)))
	fmt.Printf("%x", h)
	return h
}

type Solution struct {
	Item
	id int64
}

func NewSolution(id int64) Solution {
	sol := Solution{}
	sol.id = id
	if sol.IsNull() {
		fmt.Fprintln(os.Stderr, "warning: solution", id, "not found")
	}
	return sol
}

func (sol Solution) IsNull() bool {
	var val int
	sol.GetValue("count(*)", &val)
	return val <= 0
}

func (sol *Solution) GetHash() [16]byte {
	return md5.Sum([]byte(fmt.Sprintf("%v", sol)))
}

func (sol *Solution) SaveChanges() error {
	return errors.New("error: not implemented")
}

// GetValue gets a value from the database
func (sol *Solution) GetValue(name string, value interface{}) {
	db := currentProject.GetData()
	defer db.Close()
	err := db.GetItemValue(sol.GetId(), "Solutions", name, value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "database error:", err)
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

// SetValue sets a value to the database
func (sol *Solution) SetValue(name string, value interface{}) {
	db := currentProject.GetData()
	defer db.Close()
	db.SetItemValue(sol.GetId(), "Solutions", name, value)
}

// GetRationale gets the rationale property of the solution
func (sol *Solution) GetRationale() string {
	return sol.GetValueString("rationale")
}

// GetFitCriterion of solution
func (sol *Solution) GetFitCriterion() string {
	return sol.GetValueString("fitCriterion")
}

// GetId gets the row ID in the database
func (sol Solution) GetId() int64 {
	return sol.GetValueInt64("id")
}

// GetUid gets the row Uid in the database
func (sol Solution) GetUid() int64 {
	return sol.GetValueInt64("uid")
}

// SetUid sets the Uid in the database
func (sol Solution) SetUid(uid int64) {
	sol.SetValue("uid", uid)
}

// GetVersion of Solution
func (sol Solution) GetVersion() int {
	return sol.GetValueInt("version")
}

// GetShown gets the root as hidden or shown
func (sol Solution) GetShown() bool {
	var val bool
	sol.GetValue("shown", &val)
	return val
}

// SetShown sets the root as hidden or shown
func (sol Solution) SetShown(shown bool) {
	sol.SetValue("shown", shown)
}

// GetDescription gets the description from the database
func (sol Solution) GetDescription() string {
	return sol.GetValueString("description")
}

// AddChild adds child to solution
func (sol Solution) AddChild(child Item) {

}

func (sol Solution) Pos() (int, int) {
	var x, y int
	sol.GetValue("x", &x)
	sol.GetValue("y", &y)
	return x, y
}

func (sol Solution) SetPos(x, y int) {
	sol.SetValue("x", x)
	sol.SetValue("y", y)
}

func (sol Solution) Size() (int, int) {
	var width, height int
	sol.GetValue("width", &width)
	sol.GetValue("height", &height)
	return width, height
}

func (sol Solution) SetSize(w, h int) {
	sol.SetValue("width", w)
	sol.SetValue("height", h)
}

func (sol Solution) Parent() (parentID int64, parentType ItemType, found bool) {
	return sol.GetValueInt64("parent"), ItemType(sol.GetValueInt("parentType")), !sol.IsPropertyNull("parent")
}

func (sol Solution) IsPropertyNull(columnName string) bool {
	db := currentProject.GetData()
	defer db.Close()
	return db.IsPropertyNull("Solutions", columnName, sol.id)
}