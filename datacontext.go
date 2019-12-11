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

// DataContext is the context to the database
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

// UpdateItem updates an already existing item or creates a new one
func (data *DataContext) UpdateItem(item Item, projectVersion int) error {
	// Try to cast item as requirement
	req, isReq := item.(Requirement)
	// Get requirement/solution as string
	typeName := "s"
	if isReq {
		typeName = "r"
	}
	// Create if it doesn't exist
	if item.GetUid() == 0 {
		item.SetUid(data.GetItemUid())
		if isReq {
			// Add requirement
			stmt, err := data.Database.Prepare("insert into Requirements (uid, description, rationale, fitCriterion) values (?, ?, ?, ?)")
			_, err = stmt.Exec(item.GetUid(), item.GetDescription(), req.GetRationale(), req.GetFitCriterion())
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
			return fmt.Errorf("error: failed to get item id: %v", err)
		}
		defer stmt.Close()
		var itemID int
		err = row.Scan(&itemID)
		if err != nil || itemID == 0 {
			return fmt.Errorf("error: failed to get item id: %v", err)
		}

		stmt, err = data.Database.Prepare("insert into ItemVersions (version, item, type) values (?, ?, ?)")
		if _, err = stmt.Exec(item.GetVersion(), itemID, typeName); err != nil {

		if err != nil {
			return fmt.Errorf("error: failed to get create item version: %v", err)
		}
		defer stmt.Close()
	}

	// Rest is not completed yet
	return nil
}

//Discussion needed

//AddChild adding children to the tree
/*
func (data *DataContext) AddChild(item Item, child int, root int) {
	query := sql.NewQSqlQuery2("", data.Database)

	query.Prepare("update :tableName set parent = :parentId where id = :itemId")
	query.BindValue(":tableName", core.NewQVariant1(item.AddChild()), sql.QSql__In)
	query.BindValue(":itemId", core.NewQVariant1(item.GetChildren()), sql.QSql__In)
}

//RemoveChild delete the child (((????????? not done yet)))
func (data *DataContext) RemoveChild(item Item, child int, root int) {
	query := sql.NewQSqlQuery2("", data.Database)

	query.Prepare("update :tableName set parent = nil where id = :itemId")
	query.BindValue(":tableName", core.NewQVariant1(item.GetChildren()), sql.QSql__In)
	query.BindValue(":itemId", core.NewQVariant1(item.GetId()), sql.QSql__In)

}
*/

// UidExists checking if the specified uid is already taken
func UidExists(db *sql.DB, uid int64) bool {
	// Prepare union query
	stmt, err := db.Prepare("select count(*) from (select uid from Requirements union select uid from Solutions) where uid = ?")
	// Bind uid
	var count int
	stmt.QueryRow(uid).Scan(&count)
	// Check for error
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: failed to get uid:", err)
	}
	// If count is above 0, row is found
	return count > 0
}

//GetItemUid If the item uid dosn't exits we will get a uid
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
