import sys
from pathlib import Path
from PySide2.QtSql import QSqlDatabase, QSqlQuery
from PySide2.QtCore import QFileInfo, QJsonDocument, QRandomGenerator, QByteArray

class DataContext:
	def create(self, projectName):
		print("creating data context")
		# Load file and try to parse
		file = open("../json/tables.json", "r")
		json = QJsonDocument.fromJson(QByteArray(str.encode(file.read())))
		if json.isNull():
			print("error: failed to parse json with table information", sys.stderr)
			return False
		# Prepare query
		query = QSqlQuery(self.database)
		# Temporary list for QJSonArray
		list = []
		# Loop though all entries in the JSON
		tables = json.object()
		for key in tables.keys():
			# Clear from previous
			list.clear()
			# Transfer to list
			for row in tables[key]:
				list.append(row)
			# Execute query
			if not query.exec_("create table {0} ({1})".format(key, ", ".join(list))):
				print("database error: failed to create {0} table".format(key), sys.stderr)
				return False
		# Insert Info table
		query.prepare("insert into Info (name) values (:name)")
		query.bindValue(":name", projectName)
		if not query.exec_():
			print("database error: failed to insert into Info table", sys.stderr)
			return False
		return True

	def uidExists(self, query, uid):
		# Prepare query
		query.prepare("select count(*) from (select uid from Requirements union select uid from Solutions) where uid = :uid")
		# Bind value
		query.bindValue(":uid", uid)
		# Execute query
		query.exec_()
		query.first()
		# If above 0, it exists
		return query.value(0) > 0

	def getItemUid(self):
		def randId():
			return QRandomGenerator.global_().generate64()
		# Id to return later
		id = randId()
		# Prepare query
		query = QSqlQuery(self.database)
		# Generate new values until it doesn't exist
		while self.uidExists(query, id):
			id = randId()
		# Return UID
		return id

	def updateItem(self, item, projectVersion):
		itemType = "solution"
		if isinstance(item, Requirement):
			itemType = "requirement"
		query = QSqlQuery(self.database)
		if item.uid == 0:
			itemType.uid = self.getItemUid()
			itemId = 0
			if itemType == "requirement":
				query.prepare("insert into Requirements (uid, description, rationale, fitCriterion) values (:uid, :description, :rationale, :fitCriterion)")
				query.bindValue(":uid", item.uid)
				query.bindValue(":description", item.description)
				query.bindValue(":rationale", item.rationale)
				query.bindValue(":fitCriterion", item.fitCriterion)
			else:
				query.prepare("insert into Solutions (uid, description) values (:uid, :description)")
				query.bindValue(":uid", item.uid)
				query.bindValue(":description", item.description)
			if not query.exec_():
				print("error: failed to create item: {0}".format(self.database.lastError().text()))
				return False

			query.prepare("select id from Requirements where uid = :uid")
			query.bindValue(":uid", item.uid)
			if not query.exec_() or not query.first():
				print("error: failed to get item id: {0}".format(database.lastError().text()))
				return False
			itemId = query.value(0)

			query.prepare("insert into ItemVersions (version, item, type) values (:version, :item, :type)")
			query.bindValue(":version", projectVersion)
			query.bindValue(":item", itemId)
			query.bindValue(":type", itemType)
			if not query.exec_():
				print("error: failed to create item version", sys.stderr)
				return False
			return True
		return False
		
	def __init__(self, path):
		# Check if SQLite is available
		if QSqlDatabase.isDriverAvailable("SQLITE"):
			print("error: sqlite is not available", sys.stderr)
			return
		# Create database as SQLite database
		self.database = QSqlDatabase.addDatabase("QSQLITE")
		# Check if file already existed
		filePath = Path(path)
		fileExists = filePath.is_file()
		# Open file
		self.database.setDatabaseName(path)
		self.database.open()
		# Create database if it doesn't exist
		if not fileExists and not self.create(filePath.stem):
			print("error: failed to create initial database", sys.stderr)