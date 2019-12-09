package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/sql"
)

type DataContext struct {
	Database *sql.QSqlDatabase
}

func NewDataContext(path string) *DataContext {
	data := new(DataContext)

	// Check if SQLite is available
	if !sql.QSqlDatabase_IsDriverAvailable("QSQLITE") {
		fmt.Fprintln(os.Stderr, "error: sqlite is not available")
		return nil
	}

	// Check beforehand if file exists
	fileExists := core.QFile_Exists(path)

	// Create SQLite database
	data.Database = sql.QSqlDatabase_AddDatabase("QSQLITE", "")
	data.Database.SetDatabaseName(path)

	// Open file
	if !data.Database.Open() {
		fmt.Fprintln(os.Stderr, "error: failed to open database")
	}

	// Create if it didn't exist
	if !fileExists {
		fileName := core.NewQFileInfo3(path).FileName()
		if strings.Contains(fileName, ".") {
			fileName = strings.TrimLeft(fileName, ".")
		}
		if !data.Create(fileName) {
			fmt.Fprintln(os.Stderr, "error: failed to create initial database")
			return nil
		}
	}

	return data
}

func (data *DataContext) Create(projectName string) bool {
	// Prepare query
	query := sql.NewQSqlQuery2("", data.Database)

	// Loop through table data
	for key, value := range tableData {
		// Execute query
		if !query.Exec(fmt.Sprintf("create table %s (%s)", key, strings.Join(value, ", "))) {
			fmt.Fprintln(os.Stderr, "database error: failed to create", key, "table")
			return false
		}
	}

	// Insert info table
	query.Prepare("insert into Info (name) values (:name)")
	query.BindValue(":name", core.NewQVariant12(projectName), sql.QSql__In)
	if !query.Exec2() {
		fmt.Fprintln(os.Stderr, "database error: failed to insert into Info table")
		return false
	}

	return true
}

func UidExists(query *sql.QSqlQuery, uid int64) bool {
	// Prepare union query
	query.Prepare("select count(*) from (select uid from Requirements union select uid from Solutions) where uid = :uid")
	// Bind uid
	query.BindValue(":uid", core.NewQVariant7(uid), sql.QSql__In)
	// Execute and get result
	query.Exec2()
	query.First()
	// If count is above 0, row is found
	return query.Value(0).ToInt(nil) > 0;
}

func (data *DataContext) GetItemUid() int64 {
	// Prepare query
	query := sql.NewQSqlQuery3(data.Database)
	// Generate initial id
	id := int64(core.QRandomGenerator_Global().Generate64())
	// Keep generating until unique
	for UidExists(query, id) {
		id = int64(core.QRandomGenerator_Global().Generate64())
	}
	// Return newly generated value
	return id
}