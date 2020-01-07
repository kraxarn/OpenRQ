package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
	_, err := gz.Write(data)
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

func JSONObjectToItem(data map[string]interface{}, db *DataContext) (Item, error) {
	_, ok := data["Rationale"]
	uid, err := strconv.ParseInt(data["ID"].(string), 16, 64)
	if err != nil {
		return nil, err
	}
	if ok {
		// Has rationale, probably requirement
		id, err := db.AddRequirement(data["Description"].(string), data["Rationale"].(string), data["FitCriterion"].(string), uid)
		if err != nil {
			return nil, err
		}
		return NewRequirement(id), nil
	} else {
		// Hasn't rationale, probably solution
		id, err := db.AddSolution(data["Description"].(string), uid)
		if err != nil {
			return nil, err
		}
		return NewSolution(id), nil
	}
}

func ParseJSON(parent Item, db *DataContext, tree map[string]interface{}) error {
	// Add to DB
	item, err := JSONObjectToItem(tree, db)
	if err != nil {
		return err
	}
	// Set parent if needed
	if parent != nil {
		item.SetParent(parent)
	}
	// Set position and size
	pos := tree["Pos"].([]interface{})
	item.SetPos(int(pos[0].(float64)), int(pos[1].(float64)))
	// Do the same for children
	data, ok := tree["Children"].([]interface{})
	if ok {
		for _, child := range data {
			dataMap := child.(map[string]interface{})
			if err = ParseJSON(item, db, dataMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func NewJSONProject(path string) (*Project, error) {
	// Try to read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Try to parse json
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	// Check if valid project
	projectName, found := jsonData["ProjectName"]
	if !found {
		return nil, fmt.Errorf("not a valid project")
	}
	// Check if destination file already exists
	newPath := path[0:len(path)-4] + "orq"
	_, err = os.Stat(newPath)
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("file with name \"%v\" already exists", newPath)
	}
	// Create new data context to transfer to
	db := NewDataContext(newPath)
	// Set project name
	db.SetProjectName(projectName.(string))
	// Parse each root
	for _, root := range jsonData["Tree"].([]interface{}) {
		if err = ParseJSON(nil, db, root.(map[string]interface{})); err != nil {
			return nil, err
		}
	}
	// Everything is hopefully fine
	if err := db.Close(); err != nil {
		fmt.Println("failed to close temporary database from json:", err)
	}
	return NewProject(newPath), nil
}