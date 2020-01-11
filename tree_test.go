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
		t.Error("unexpected item count, expected 0, but found", count)
	}
	// Add a new requirement
	_, err = db.AddEmptyRequirement()
	if err != nil {
		t.Error("failed to add requirement 1:", err)
	}
	// Close database
	if err = db.Close(); err != nil {
		t.Error("failed to close database connection:", err)
	}
}