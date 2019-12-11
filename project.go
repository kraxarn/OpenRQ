package main

import (
	"strings"
)

type Project struct {
	Open bool
	path string
}

func NewProject(path string) *Project {
	p := new(Project)
	p.Open = true
	p.path = path
	if !strings.HasSuffix(path, ".orq") {
		path += ".orq"
	}
	NewDataContext(path).Close()
	return p
}

func (proj *Project) GetData() *DataContext {
	return NewDataContext(proj.path)
}
