package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// ItemType enum
type ItemType int8

const (
	TypeSolution    ItemType = 0
	TypeRequirement ItemType = 1
)

// DataContext holding the connection to the database
type DataContext struct {
	Database *sql.DB
}

// NewDataContext creates a new database
func NewDataContext(path string) *DataContext {
	data := new(DataContext)

	// Check beforehand if file exists
	_, err := os.Stat(path)
	fileExists := !os.IsNotExist(err)

	// Create SQLite database
	data.Database, err = sql.Open("sqlite3", path)

	// Open file
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to open database:", err)
		return nil
	}

	// Create if it didn't exist
	if !fileExists {
		fileName := filepath.Base(path)
		if strings.Contains(fileName, ".") {
			fileName = fileName[0:strings.LastIndex(fileName, ".")]
		}
		if err := data.Create(fileName); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to create database:", err)
			return nil
		}
	}

	return data
}

// Close closes the connection to the database
func (data *DataContext) Close() error {
	return data.Database.Close()
}

// Create creates a new, empty database
func (data *DataContext) Create(projectName string) error {
	// Loop through table data
	for key, value := range tableData {
		// Execute query
		_, err := data.Database.Exec(fmt.Sprintf("create table %s (%s)", key, strings.Join(value, ", ")))
		// Check for errors
		if err != nil {
			return fmt.Errorf("error: failed to create %v table", key)
		}
	}
	// Insert info table
	_, err := data.Database.Exec("insert into Info (name) values (?)", projectName)
	return err
}

func (data *DataContext) ProjectName() string {
	var name string
	row := data.Database.QueryRow("select name from Info")
	if err := row.Scan(&name); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to get project name:", err)
	}
	return name
}

func (data *DataContext) SetProjectName(name string) {
	if _, err := data.Database.Exec("update Info set name = ?", name); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set project name:", err)
	}
}

func (data *DataContext) AddEmptyRequirement() (int64, error) {
	return data.AddRequirement("", "", "")
}

// AddRequirement adds a requirement to the current database
func (data *DataContext) AddRequirement(description, rationale, fitCriterion string) (int64, error) {
	// Generate a new random UID
	reqUID := data.ItemUID()
	// Try to add the requirement to the database
	if _, err := data.Database.Exec(
		"insert into Requirements (uid, description, rationale, fitCriterion) values (?, ?, ?, ?)",
		reqUID, description, rationale, fitCriterion); err != nil {
		return 0, err
	}
	// Get item ID
	var id int64
	if err := data.Database.QueryRow(
		"select _rowid_ from Requirements where uid = ?", reqUID).Scan(&id); err != nil {
		return 0, err
	}
	// Try to version it and return the result of it
	return id, data.AddItemVersion(reqUID, TypeRequirement)
}

// AddSolution adds a solution to the database
func (data *DataContext) AddSolution(description string) (int64, error) {
	// Generate a new random UID
	solUID := data.ItemUID()
	// Try to add the solution to the database
	if _, err := data.Database.Exec(
		"insert into Solutions (uid, description) values (?, ?)",
		solUID, description); err != nil {
		return 0, err
	}
	// Try to version it and return the result of it
	return solUID, data.AddItemVersion(solUID, TypeSolution)
}

// AddItemVersion versions an item
func (data *DataContext) AddItemVersion(itemUID int64, itemType ItemType) error {
	// Find item ID
	var itemID int
	row := data.Database.QueryRow(
		fmt.Sprintf("select _rowid_ from %v where uid = ?", GetItemTableName(itemType)), itemUID)
	if err := row.Scan(&itemID); err != nil {
		return err
	}
	// Insert into item versions
	_, err := data.Database.Exec(
		"insert into ItemVersions (version, item, type) values (1, ?, ?)",
		itemID, itemType)
	return err
}

// RemoveItem removes the item from the current version
func (data *DataContext) RemoveItem(itemID int64) error {
	// TODO: Make difference between solution and requirements
	// Execute SQL
	_, err := data.Database.Exec("delete from Requirements where _rowid_ = ?", itemID)
	return err
}

// UpdateItem updates the content of an item
func (data *DataContext) UpdateItem(item Item, projectVersion int) error {
	// Try to cast item as requirement
	req, isReq := item.(Requirement)
	// Error checking
	if item.ID() == 0 || item.UID() == 0 {
		return fmt.Errorf("failed to update item, does not exist")
	}

	itemUID := data.ItemUID()
	if isReq {
		// Add requirement
		_, err := data.Database.Exec("insert into Requirements (uid) values (?)", itemUID)
		if err != nil {
			return fmt.Errorf("failed to insert requirement: %v", err)
		}
	} else {
		// Add solution
		_, err := data.Database.Exec("insert into Solutions (uid, description) values (?, ?)")
		if err != nil {
			return fmt.Errorf("failed to insert solution: %v", err)
		}
	}

	// Find item id
	row := data.Database.QueryRow("select _rowid_ from Requirements where uid = ?", item.UID())
	if err := row.Scan(&req.id); err != nil || req.id == 0 {
		return fmt.Errorf("failed to get item id: %v", err)
	}

	// Rest is not completed yet
	return nil
}

// GetItemType gets what struct the generic item interface is
func GetItemType(item Item) ItemType {
	switch item.(type) {
	case Requirement:
		return TypeRequirement
	default:
		return TypeSolution
	}
}

// GetItemTableName gets the table name in the database for an item
func GetItemTableName(itemType ItemType) string {
	switch itemType {
	case TypeRequirement:
		return "Requirements"
	default:
		return "Solutions"
	}
}

// GetAllItems gets all requirements and solutions stored in the database
func (data *DataContext) Items() (items map[Item]string, err error) {
	// Connect to database
	db := currentProject.Data()
	defer db.Close()
	// Crate slice of items
	items = make(map[Item]string)
	// Get all requirements
	rows, err := data.Database.Query("select _rowid_, description from Requirements")
	if err != nil {
		return items, fmt.Errorf("failed to get requirements: %v", err)
	}
	var itemID int64
	var description string
	for rows.Next() {
		if err = rows.Scan(&itemID, &description); err == nil {
			items[NewRequirement(itemID)] = description
		}
	}
	// Get all solutions
	rows, err = data.Database.Query("select _rowid_, description from Solutions")
	if err != nil {
		return items, fmt.Errorf("failed to get solutions: %v", err)
	}
	for rows.Next() {
		if err := rows.Scan(&itemID, &description); err == nil {
			items[NewSolution(itemID)] = description
		}
	}
	return items, nil
}


func (data *DataContext) Links() (items map[Item]Item, err error) {
	// Connect to database
	db := currentProject.Data()
	defer db.Close()
	// Create map
	items = make(map[Item]Item)
	// Get all requirements
	// TODO: Only child and parent as requirement
	rows, err := db.Database.Query("select _rowid_, parent from Requirements where parent is not null")
	if err != nil {
		return items, fmt.Errorf("failed to get requirement links: %v", err)
	}
	defer rows.Close()
	// Get requirement links
	var itemID, parentID int64
	for rows.Next() {
		if err = rows.Scan(&itemID, &parentID); err == nil {
			items[NewRequirement(itemID)] = NewRequirement(parentID)
		} else {
			return items, fmt.Errorf("failed to get requirement %v link: %v", itemID, err)
		}
	}
	return items, nil
}

// GetItemValue gets a value from the specified column in the database
// (value is assumed to be a pointer)
func (data *DataContext) GetItemValue(itemID int64, tableName, name string, value interface{}) error {
	row := data.Database.QueryRow(fmt.Sprintf("select %v from %v where _rowid_ = ?", name, tableName), itemID)
	// Execute and return it
	if err := row.Scan(value); err != nil {
		return fmt.Errorf("failed to get property %v from item %v: %v", name, itemID, err)
	}
	return nil
}

func (data *DataContext) IsItemPropertyNull(itemID int64, tableName, columnName string) bool {
	var count int
	row := data.Database.QueryRow(
		fmt.Sprintf("select count(*) from %v where _rowid_ = ? and %v is null", tableName, columnName), itemID)
	if err := row.Scan(&count); err != nil {
		fmt.Println(err)
	}
	return count > 0
}

// AddItemChild creates a link between parent and child
func (data *DataContext) AddItemChild(parent, child Item) error {
	// Get what table name child has
	childTable := GetItemTableName(GetItemType(child))
	// Execute update
	_, err := data.Database.Exec(
		fmt.Sprintf("update %v set parent = ?, parentType = ? where _rowid_ = ?", childTable),
		parent.ID(), GetItemType(parent), child.ID())
	return err
}

// RemoveItemParent removes the link between parent and child
func (data *DataContext) RemoveChildrenLinks(parent Item) error {
	// Execute update
	// TODO: Assume requirement
	_, err := data.Database.Exec(
		"update Requirements set parent = null, parentType = null where parent = ?", parent.ID())
	return err
}

// SetItemValue updates a value in the database
func (data *DataContext) SetItemValue(itemID int64, tableName, name string, value interface{}) {
	_, err := data.Database.Exec(
		fmt.Sprintf("update %v set %v = ? where _rowid_ = ?", tableName, name), value, itemID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set property", name, "in requirement:", err)
	}
}

// UidExists checking if the specified uid is already taken
func (data *DataContext) UIDExists(uid int64) bool {
	// Execute query
	row := data.Database.QueryRow(
		"select count(*) from (select uid from Requirements union select uid from Solutions) where uid = ?", uid)
	// Get value from row
	var count int
	row.Scan(&count)
	return count > 0
}

// GetItemUid If the item uid doesn't exits we will get a uid
func (data *DataContext) ItemUID() int64 {
	// Generate initial id
	id := int64(rand.Uint64())
	// Keep generating until unique
	for data.UIDExists(id) {
		id = int64(rand.Uint64())
	}
	// Return newly generated value
	return id
}
