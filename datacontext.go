package main

import (
	"fmt"
	"math/rand"
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

func (data *DataContext) UpdateItem(item Item, projectVersion int) error {
	// Try to cast item as requirement
	req, isReq := item.(Requirement)
	// Get requirement/solution as string
	typeName := "s"
	if isReq {
		typeName = "r"
	}
	// Prepare query
	query := sql.NewQSqlQuery3(data.Database)
	// Create if it doesn't exist
	if item.GetUid() == 0 {
		item.SetUid(data.GetItemUid())
		if isReq {
			// Add requirement
			query.Prepare("insert into Requirements (uid, description, rationale, fitCriterion) values (:uid, :description, :rationale, :fitCriterion)")
			query.BindValue(":uid", core.NewQVariant7(item.GetUid()), sql.QSql__In)
			query.BindValue(":description", core.NewQVariant12(item.GetDescription()), sql.QSql__In)
			query.BindValue(":rationale", core.NewQVariant12(req.GetRationale()), sql.QSql__In)
			query.BindValue(":fitCriterion", core.NewQVariant12(req.GetFitCriterion()), sql.QSql__In)
		} else {
			// Add solution
			query.Prepare("insert into Solutions (uid, description) values (:uid, :description)")
			query.BindValue(":uid", core.NewQVariant7(item.GetUid()), sql.QSql__In)
			query.BindValue(":description", core.NewQVariant12(item.GetDescription()), sql.QSql__In)
		}

		// Try to add item
		if !query.Exec2() {
			return fmt.Errorf("failed to create item: %v", data.Database.LastError().Text())
		}

		// Find item id
		query.Prepare("select id from Requirements where uid = :uid")
		query.BindValue(":uid", core.NewQVariant7(item.GetUid()), sql.QSql__In)
		if !query.Exec2() {
			return fmt.Errorf("failed to get item id: %v", data.Database.LastError().Text())
		}

		query.Prepare("insert into ItemVersions (version, item, type) values (:version, :item, :type)")
		query.BindValue(":version", core.NewQVariant5(item.GetVersion()), sql.QSql__In)
		query.BindValue(":item", core.NewQVariant5(item.GetId()), sql.QSql__In)
		query.BindValue(":type", core.NewQVariant12(typeName), sql.QSql__In)

		if !query.Exec2() {
			return fmt.Errorf("failed to get create item version: %v", data.Database.LastError().Text())
		}
	}

	// Rest is not completed yet
	return nil
}

//Discussion needed

//AddChild adding children to the tree
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

//UidExists checking if the specified uid is already taken
func UidExists(query *sql.QSqlQuery, uid int64) bool {
	// Prepare union query
	query.Prepare("select count(*) from (select uid from Requirements union select uid from Solutions) where uid = :uid")
	// Bind uid
	query.BindValue(":uid", core.NewQVariant7(uid), sql.QSql__In)
	// Execute and get result
	query.Exec2()
	query.First()
	// If count is above 0, row is found
	return query.Value(0).ToInt(nil) > 0
}

//GetItemUid If the item uid dosn't exits we will get a uid
func (data *DataContext) GetItemUid() int64 {
	// Prepare query
	query := sql.NewQSqlQuery3(data.Database)
	// Generate initial id
	id := int64(rand.Uint64())
	// Keep generating until unique
	for UidExists(query, id) {
		id = int64(rand.Uint64())
	}
	// Return newly generated value
	return id
}
