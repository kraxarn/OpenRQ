package main

import (
	"strings"
)

type Project struct {
	Open bool
	data *DataContext
}

func NewProject(path string) *Project {
	p := new(Project)
	p.Open = true
	if !strings.HasSuffix(path, ".orq") {
		path += ".orq"
	}
	p.data = NewDataContext(path)
	return p
}
