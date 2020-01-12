package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTree(t *testing.T) {
	// Get temporary directory
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Error("failed to get temporary directory:", err)
		return
	}
	// Check if a file with the specified name already exists
	projectPath := fmt.Sprintf("%v/openrq_test.orq", tempDir)
	if _, err = os.Stat(projectPath); os.IsExist(err) {
		t.Error("failed to create temporary project: file exists")
		return
	}
	// Create a temporary project there
	project := NewProject(projectPath)
	if project == nil {
		t.Error("failed to create new project")
		return
	}
	// Create database connection
	db := project.Data()
	// Check so that we have no items
	if count := len(Roots()); count != 0 {
		t.Error("unexpected item count, expected 0, but got", count)
	}
	// Add a new requirement
	req1id, err := db.AddEmptyRequirement()
	if err != nil {
		t.Error("failed to add requirement 1:", err)
	}
	req1 := NewRequirement(req1id)
	// Make sure the requirement description is empty
	if descLen := len(req1.Description()); descLen != 0 {
		t.Error("unexpected requirement 1 description length, expected 0, but got", descLen)
	}
	// Try to set a description
	req1desc := "requirement 1 description"
	req1.SetDescription(req1desc)
	if req1.Description() != req1desc {
		t.Errorf("unexpected requirement 1 description, expected \"%v\", but got \"%v\"",
			req1desc, req1.Description())
	}
	// Make sure the requirement doesn't have any children
	if count := len(req1.Children()); count != 0 {
		t.Error("unexpected requirement 1 children count, expected 0, but got", count)
	}
	// Create a new solution
	sol1id, err := db.AddEmptySolution()
	if err != nil {
		t.Error("failed to add solution 1:", err)
	}
	sol1 := NewSolution(sol1id)
	// Make sure solution 1 has no parent
	if sol1.Parent() != nil {
		t.Error("unexpected solution 1 parent, expected nil, but got", sol1.Parent().ToString())
	}
	// Link requirement 1 to solution 1
	if err = db.AddItemChild(req1, sol1); err != nil {
		t.Error("failed to create link between requirement 1 and solution 1:", err)
	}
	// Make sure solution 1 has requirement 1 as parent
	if sol1.Parent() != req1 {
		t.Error("unexpected solution 1 parent, expected requirement 1, but got", sol1.Parent().ToString())
	}
	// Close database
	if err = db.Close(); err != nil {
		t.Error("failed to close database connection:", err)
	}
}