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

func NewCompressedProject(path string) (*Project, error) {
	// Name of final output file
	newPath := path[0:len(path)-1]
	// Check if decompressed file already exists
	_, err := os.Stat(newPath)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("file with name \"%v\" already exists", newPath)
	}
	// Get file permissions of original file for later
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	// Load compressed file
	compressed, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Decompress it
	decompressed, err := Decompress(compressed)
	if err != nil {
		return nil, err
	}
	// Create new decompressed copy
	if err := ioutil.WriteFile(newPath, decompressed, fileInfo.Mode()); err != nil {
		return nil, err
	}
	// Remove old project
	if err := os.Remove(path); err != nil {
		return nil, err
	}
	// Everything seems fine, load newly created project
	return NewProject(newPath), nil
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