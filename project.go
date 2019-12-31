package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
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

func Compress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	count, err := gz.Write(data)
	_ = gz.Close()
	fmt.Println("wrote", count, "/", len(buffer.Bytes()), "bytes")
	return buffer.Bytes(), err
}

func Decompress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	_ = gz.Close()
	return ioutil.ReadAll(gz)
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
	// GZip
	if strings.HasSuffix(path, ".orqz") {
		temp, err := Compress(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed to compress data:", err)
		} else {
			file = temp
		}
	}
	// Copy data to location
	return ioutil.WriteFile(path, file, fileInfo.Mode())
}