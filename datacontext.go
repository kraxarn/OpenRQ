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
	stmt, err := data.Database.Prepare("insert into Info (name) values (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(projectName)
	return err
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
	if err := data.Database.QueryRow("select _rowid_ from Requirements where uid = ?", reqUID).Scan(&id); err != nil {
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
	stmt, err := data.Database.Prepare(fmt.Sprintf("select _rowid_ from %v where uid = ?", GetItemTableName(itemType)))
	if err != nil {
		return fmt.Errorf("failed to get item id: %v", err)
	}
	stmt.QueryRow(itemUID).Scan(&itemID)
	// Insert into item versions
	_, err = data.Database.Exec(
		"insert into ItemVersions (version, item, type) values (1, ?, ?)",
		itemID, itemType)
	return err
}

// UpdateItemVersion increases the version of an item
// TODO: Not yet implemented as version control isn't implemented
func (data *DataContext) UpdateItemVersion() error {
	return nil
}

// RemoveItem removes the item from the current version
func (data *DataContext) RemoveItem(itemType ItemType, itemID int) error {
	// Execute SQL
	_, err := data.Database.Exec(
		"delete from ? where _rowid_ = ?",
		GetItemTableName(itemType), itemID)
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
		stmt, err := data.Database.Prepare("insert into Requirements (uid) values (?)")
		if err != nil {
			return fmt.Errorf("failed to insert requirement: %v", err)
		}
		_, err = stmt.Exec(itemUID)
		defer stmt.Close()
	} else {
		// Add solution
		stmt, err := data.Database.Prepare("insert into Solutions (uid, description) values (?, ?)")
		if _, err = stmt.Exec(item.UID(), item.Description()); err != nil {
			return fmt.Errorf("failed to insert solution: %v", err)
		}
		defer stmt.Close()
	}

	// Find item id
	stmt, err := data.Database.Prepare("select _rowid_ from Requirements where uid = ?")
	if err != nil {
		return fmt.Errorf("failed to get item id: %v", err)
	}
	row := stmt.QueryRow(item.UID())
	defer stmt.Close()
	err = row.Scan(&req.id)
	if err != nil || req.id == 0 {
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
func (data *DataContext) Items() ([]Item, error) {
	// Connect to database
	db := currentProject.Data()
	defer db.Close()
	// Temporary slice
	items := make([]Item, 0)
	// Get all requirements
	stmt, err := data.Database.Prepare("select _rowid_ from Requirements")
	if err != nil {
		return items, fmt.Errorf("failed to get requirements: %v", err)
	}
	rows, _ := stmt.Query()
	var itemID int64
	for rows.Next() {
		if err = rows.Scan(&itemID); err == nil {
			items = append(items, NewRequirement(itemID))
		}
	}
	// Get all solutions
	stmt, err = data.Database.Prepare("select _rowid_ from Solutions")
	if err != nil {
		return items, fmt.Errorf("failed to get solutions: %v", err)
	}
	rows, _ = stmt.Query()
	for rows.Next() {
		if err := rows.Scan(&itemID); err == nil {
			items = append(items, NewRequirement(itemID))
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
			items[NewRequirement(parentID)] = NewRequirement(itemID)
		} else {
			return items, fmt.Errorf("failed to get requirement %v link: %v", itemID, err)
		}
	}
	return items, nil
}

// GetItemChildren gets all children of a specific item
func (data *DataContext) GetItemChildren(itemID int) {
	stmt, err := data.Database.Prepare("select parent from (select parent from Requirements union select parent from Solutions) where parent = ? and parentType = ?")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to prepare get of item", itemID, ":", err)
	}
	if _, err := stmt.Query(itemID); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to get item", itemID, ":", err)
	}
}

// GetItemValue gets a value from the specified column in the database
// (value is assumed to be a pointer)
func (data *DataContext) GetItemValue(itemID int64, tableName, name string, value interface{}) error {
	// Prepare query
	stmt, err := data.Database.Prepare(fmt.Sprintf("select %v from %v where _rowid_ = ?", name, tableName))
	if err != nil {
		return fmt.Errorf("failed to get property %v from item %v: %v", name, itemID, err)
	}
	// Execute and return it
	if err := stmt.QueryRow(itemID).Scan(value); err != nil {
		return fmt.Errorf("failed to get property %v from item %v: %v", name, itemID, err)
	}
	return nil
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
func (data *DataContext) RemoveItemParent(child Item) error {
	// Get what table name child has
	childTable := "Solutions"
	if GetItemType(child) == TypeRequirement {
		childTable = "Requirements"
	}
	// Execute update
	_, err := data.Database.Exec(
		"update ? set parent = null and parentType = null where _rowid_ = ?",
		childTable, child.ID())
	return err
}

// SetItemValue updates a value in the database
func (data *DataContext) SetItemValue(itemID int64, tableName, name string, value interface{}) {
	// Prepare query
	stmt, err := data.Database.Prepare(fmt.Sprintf("update %v set %v = ? where _rowid_ = ?", tableName, name))
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set property", name, "in requirement:", err)
		return
	}
	// Execute and return it
	if _, err := stmt.Exec(value, itemID); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set property", name, "in requirement:", err)
	}
}

func (data *DataContext) IsPropertyNull(tableName, columnName string, id int64) bool {
	// Prepare query
	stmt, err := data.Database.Prepare(
		fmt.Sprintf("select count(*) from %v where %v is null and _rowid_ = %v", tableName, columnName, id))
	if err != nil {
		fmt.Println("error: failed to check for null property:", err)
		return true
	}
	// Execute and return it
	var value int
	_ = stmt.QueryRow(0).Scan(&value)
	return value >= 0
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
