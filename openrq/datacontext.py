import sys
from pathlib import Path
from PySide2.QtSql import QSqlDatabase

class DataContext:
	def __init__(self, path):
		# Check if SQLite is available
		if QSqlDatabase.isDriverAvailable("SQLITE"):
			print("error: sqlite is not available", sys.stderr)
			return
		# Create database as SQLite database
		self.database = QSqlDatabase.addDatabase("QSQLITE")
		# Check if file already existed
		filePath = Path(path)
		# Open file
		self.database.setDatabaseName(path)
		self.database.open()
		# Create database if it doesn't exist
		if not filePath.is_file() and not self.create(filePath.stem):
			print("error: failed to create initial database")
	
	def create(self, projectName):
	        file = open("../json/tables.json", "r")
