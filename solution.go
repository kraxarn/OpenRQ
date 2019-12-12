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

// GetHash gets the md5 hash of the item
func (link *Link) GetHash() [16]byte {
	h := md5.Sum([]byte(fmt.Sprintf("%v%v", link.name, link.color)))
	fmt.Printf("%x", h)
	return h
}

type Solution struct {
	id int
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

// GetValue gets a value from the database
func (sol *Solution) GetValue(name string) interface{} {
	db := currentProject.GetData()
	defer db.Close()
	return db.GetItemValue(sol.GetId(), "Solutions", name)
}

// SetValue sets a value to the database
func (sol *Solution) SetValue(name string, value interface{}) {
	db := currentProject.GetData()
	defer db.Close()
	db.SetItemValue(sol.GetId(), "Solutions", name, value)
}

// GetRationale gets the rationale property of the solution
func (sol *Solution) GetRationale() string {
	return sol.GetValue("rationale").(string)
}

// GetFitCriterion of solution
func (sol *Solution) GetFitCriterion() string {
	return sol.GetValue("fitCriterion").(string)
}

// GetId gets the row ID in the database
func (sol Solution) GetId() int {
	return sol.GetValue("id").(int)
}

// GetUid gets the row Uid in the database
func (sol Solution) GetUid() int64 {
	return sol.GetValue("uid").(int64)
}

// SetUid sets the Uid in the database
func (sol Solution) SetUid(uid int64) {
	sol.SetValue("uid", uid)
}

// GetVersion of Solution
func (sol Solution) GetVersion() int {
	return sol.GetValue("version").(int)
}

// GetShown gets the root as hidden or shown
func (sol Solution) GetShown() bool {
	return sol.GetValue("shown").(bool)
}

// SetShown sets the root as hidden or shown
func (sol Solution) SetShown(shown bool) {
	sol.SetValue("shown", shown)
}

// GetDescription gets the description from the database
func (sol Solution) GetDescription() string {
	return sol.GetValue("description").(string)
}

// AddChild adds child to solution
func (sol Solution) AddChild(child Item) {

}
