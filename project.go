package main

import (
	"io/ioutil"
	"os"
	"strings"
)

var currentProject *Project

type Project struct {
	Open bool
	path string
}

func NewProject(path string) *Project {
	currentProject = new(Project)
	currentProject.Open = true
	currentProject.path = path
	if !strings.HasSuffix(path, ".orq") {
		path += ".orq"
	}
	NewDataContext(path).Close()
	return currentProject
}

func (proj *Project) Data() *DataContext {
	return NewDataContext(proj.path)
}

func (proj *Project) Name() string {
	db := proj.Data()
	defer db.Close()
	return db.ProjectName()
}

func (proj *Project) CopyTo(path string) error {
	// Get file info to copy permissions
	fileInfo, err := os.Stat(proj.path)
	if err != nil {
		return err
	}
	// Get the current project file
	file, err := ioutil.ReadFile(proj.path)
	if err != nil {
		return err
	}
	// Copy data to location
	return ioutil.WriteFile(path, file, fileInfo.Mode())
}