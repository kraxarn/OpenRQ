package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

// ValidateLinkErrors can't be easily tested
// as they depend on the graphics view

func TestValidateLinks(t *testing.T) {
	links = make(map[Item][]*Link)
	// Create two items of different types
	req1 := NewRequirement(1)
	sol1 := NewSolution(1)
	// Create a link between the items
	link := Link{
		parent: req1,
		child:  sol1,
	}
	links[req1] = append(links[req1], &link)
	links[sol1] = append(links[sol1], &link)
	// This is valid, so link validation should be successful
	if links := len(ValidateLinks()); links != 0 {
		t.Error("unexpected validation result, expected 0 errors, but got", links)
	}
	// Try adding a new solution to the current solution
	sol2 := NewSolution(2)
	link2 := Link{
		parent: sol1,
		child:  sol2,
	}
	links[sol1] = append(links[sol1], &link2)
	links[sol2] = append(links[sol2], &link2)
	// Validation should fail now
	if links := len(ValidateLinks()); links != 1 {
		t.Error("unexpected validation result, expected 1 error, but got", links)
	}
}

func TestValidateLoops(t *testing.T) {
	links = make(map[Item][]*Link)
	// Create two items of different types
	req1 := NewRequirement(1)
	sol1 := NewSolution(1)
	// Create a link between the items
	link := Link{
		parent: req1,
		child:  sol1,
	}
	links[req1] = append(links[req1], &link)
	links[sol1] = append(links[sol1], &link)
	// This is valid, so link validation should be successful
	if links := len(ValidateLoops()); links != 0 {
		t.Error("unexpected validation result, expected 0 errors, but got", links)
	}
	// Try linking them again, but in the other direction
	link2 := Link{
		parent: sol1,
		child:  req1,
	}
	links[sol1] = append(links[sol1], &link2)
	links[req1] = append(links[req1], &link2)
	// Validation should fail now
	if links := len(ValidateLoops()); links != 2 {
		t.Error("unexpected validation result, expected 1 error, but got", links)
	}
}

func TestValidateRoots(t *testing.T) {
	// ValidateRoots depend on a view being available
	_, window := NewMainWindow()
	CreateLayout(window)
	// Reset project if any was loaded
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error(err)
	}
	NewProject(fmt.Sprintf("%v/openrq_test.orq", tmpDir))
	ReloadProject(window)
	// Create two items of different types
	req1 := NewGraphicsItem("req1", 0, 0, 0, 0, NewRequirement(1))
	sol1 := NewGraphicsItem("sol1", 0, 0, 0, 0, NewSolution(1))
	// Add to scene
	scene.AddItem(req1)
	scene.AddItem(sol1)
	// Create a link between the items
	CreateLink(req1, sol1)
	// Validation should be successful
	if links := len(ValidateRoots()); links != 0 {
		t.Error("unexpected validation result, expected 0 errors, but got", links)
	}
	// Add another child to parent
	sol2 := NewGraphicsItem("sol2", 0, 0, 0, 0, NewSolution(2))
	scene.AddItem(sol2)
	CreateLink(req1, sol2)
	// Validation should now fail
	if links := len(ValidateRoots()); links != 1 {
		t.Error("unexpected validation result, expected 1 error, but got", links)
	}
}