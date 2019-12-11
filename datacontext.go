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

type ItemType int8

const (
	TypeSolution    ItemType = 0
	TypeRequirement ItemType = 1
)

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
			fileName = strings.TrimLeft(fileName, ".")
		}
		if err := data.Create(fileName); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to create database:", err)
			return nil
		}
	}

	return data
}

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

func (data *DataContext) AddRequirement(description, rationale, fitCriterion string) error {
	reqUID := data.GetItemUid()
	if _, err := data.Database.Exec(
		"insert into Requirements (uid, description, rationale, fitCriterion) values (?, ?, ?, ?)",
		reqUID, description, rationale, fitCriterion); err != nil {
		return err
	}
	return data.AddItemVersion(reqUID, TypeRequirement)
}

func (data *DataContext) AddSolution(description string) error {
	solUID := data.GetItemUid()
	if _, err := data.Database.Exec(
		"insert into Solutions (uid, description) values (?, ?)",
		solUID, description); err != nil {
		return err
	}

	return data.AddItemVersion(solUID, TypeSolution)
}

func (data *DataContext) AddItemVersion(itemUID int64, itemType ItemType) error {
	// Find item id
	var itemID int
	stmt, err := data.Database.Prepare("select id from ? where uid = ?")
	if err != nil {
		return fmt.Errorf("failed to get item id: %v", err)
	}
	stmt.QueryRow(GetItemTableName(itemType), itemUID).Scan(&itemID)
	// Insert into item versions
	_, err = data.Database.Exec(
		"insert into ItemVersions (version, item, type) values (1, ?, ?)",
		itemID, itemType)
	return err
}

// UpdateItemVersion Implement this later
func (data *DataContext) UpdateItemVersion() error {
	return nil
}

func (data *DataContext) RemoveItem(itemType ItemType, itemID int) error {
	// Execute SQL
	_, err := data.Database.Exec(
		"delete from ? where id = ?",
		GetItemTableName(itemType), itemID)
	return err
}

func (data *DataContext) UpdateItem(item Item, projectVersion int) error {
	// Try to cast item as requirement
	req, isReq := item.(Requirement)
	// Error checking
	if item.GetId() == 0 || item.GetUid() == 0 {
		return fmt.Errorf("failed to update item, does not exist")
	}

	itemUID := data.GetItemUid()
	if isReq {
		// Add requirement
		stmt, err := data.Database.Prepare("insert into Requirements (uid) values (?)")
		_, err = stmt.Exec(itemUID)
		if err != nil {
			return fmt.Errorf("failed to insert requirement: %v", err)
		}
		defer stmt.Close()
	} else {
		// Add solution
		stmt, err := data.Database.Prepare("insert into Solutions (uid, description) values (?, ?)")
		if _, err = stmt.Exec(item.GetUid(), item.GetDescription()); err != nil {
			return fmt.Errorf("failed to insert solution: %v", err)
		}
		defer stmt.Close()
	}

	// Find item id
	stmt, err := data.Database.Prepare("select id from Requirements where uid = ?")
	row := stmt.QueryRow(item.GetUid())
	if err != nil {
		return fmt.Errorf("failed to get item id: %v", err)
	}
	defer stmt.Close()
	err = row.Scan(&req.id)
	if err != nil || req.id == 0 {
		return fmt.Errorf("failed to get item id: %v", err)
	}

	// Rest is not completed yet
	return nil
}

func GetItemType(item Item) ItemType {
	switch item.(type) {
	case Requirement:
		return TypeRequirement
	default:
		return TypeSolution
	}
}

func GetItemTableName(itemType ItemType) string {
	switch itemType {
	case TypeRequirement:
		return "Requirements"
	default:
		return "Solutions"
	}
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
func (data *DataContext) GetItemValue(itemID int, tableName, name string) interface{} {
	// Prepare query
	stmt, err := data.Database.Prepare("select ? from ? where id = ?")
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to get property", name, "from requirement:", err)
		return nil
	}
	// Execute and return it
	var value interface{}
	if err := stmt.QueryRow(itemID, tableName, name).Scan(&value); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to get property", name, "from requirement:", err)
		return nil
	}
	return value
}

func (data *DataContext) AddItemChild(parent, child Item) error {
	// Get what table name child has
	childTable := "Solutions"
	if GetItemType(child) == TypeRequirement {
		childTable = "Requirements"
	}
	// Execute update
	_, err := data.Database.Exec(
		"update ? set parent = ? and parentType = ? where id = ?",
		childTable, parent.GetId(), GetItemType(parent), child.GetId())
	return err
}

func (data *DataContext) RemoveItemParent(child Item) error {
	// Get what table name child has
	childTable := "Solutions"
	if GetItemType(child) == TypeRequirement {
		childTable = "Requirements"
	}
	// Execute update
	_, err := data.Database.Exec(
		"update ? set parent = null and parentType = null where id = ?",
		childTable, child.GetId())
	return err
}

// SetItemValue updates a value in the database
func (data *DataContext) SetItemValue(itemID int, tableName, name string, value interface{}) {
	// Prepare query
	stmt, err := data.Database.Prepare("update ? set ? = ? where id = ?")
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set property", name, "in requirement:", err)
		return
	}
	// Execute and return it
	if _, err := stmt.Exec(tableName, name, value, itemID); err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to set property", name, "in requirement:", err)
	}
}

// UidExists checking if the specified uid is already taken
func UidExists(db *sql.DB, uid int64) bool {
	// Prepare union query
	stmt, err := db.Prepare("select count(*) from (select uid from Requirements union select uid from Solutions) where uid = ?")
	defer stmt.Close()
	// Check for error
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to get uid:", err)
	}
	// Bind uid
	var count int
	stmt.QueryRow(uid).Scan(&count)
	// If count is above 0, row is found
	return count > 0
}

// GetItemUid If the item uid doesn't exits we will get a uid
func (data *DataContext) GetItemUid() int64 {
	// Generate initial id
	id := int64(rand.Uint64())
	// Keep generating until unique
	for UidExists(data.Database, id) {
		id = int64(rand.Uint64())
	}
	// Return newly generated value
	return id
}
