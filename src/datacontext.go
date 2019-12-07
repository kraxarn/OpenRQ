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

func (data *DataContext) Create(projectName string) bool {
	// Load file and try to parse
	file := core.NewQFile2("../json/tables.json")
	if !file.Open(core.QIODevice__ReadOnly | core.QIODevice__Text) {
		fmt.Fprintln(os.Stderr, "error: failed to open json")
		return false
	}
	json := core.QJsonDocument_FromJson(file.ReadAll(), nil)
	if json.IsNull() {
		fmt.Fprintln(os.Stderr, "error: failed to parse json")
		return false
	}

	// Prepare query
	query := sql.NewQSqlQuery2("", data.Database)

	// Temporary list for QJsonArray
	list := []string{}

	// Loop through JSON file
	tables := json.Object()
	for _, key := range tables.Keys() {
		// Clear from previous
		list = list[:0]

		// Transfer to list
		for _, row := range tables.Value(key).ToArray2().ToVariantList() {
			list = append(list, row.ToString())
		}

		// Execute query
		if !query.Exec(fmt.Sprintf("create table %s (%s)", key, strings.Join(list, ", "))) {
			fmt.Fprintln(os.Stderr, "database error: failed to create %s table", key)
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

func NewDataContext(path string) *DataContext {
	data := new(DataContext)

	// Check if SQLite is available
	if sql.QSqlDatabase_IsDriverAvailable("SQLITE") {
		fmt.Fprintln(os.Stderr, "error: sqlite is not available")
		return nil
	}

	// Check beforehand if file exists
	fileExists := core.QFile_Exists(path)

	// Create SQLite database
	data.Database = sql.QSqlDatabase_AddDatabase("QSQLITE", path)

	// Open file
	data.Database.Open()

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
