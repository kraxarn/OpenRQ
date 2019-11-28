#include "datacontext.h"

DataContext::DataContext(QString path)
{
	// Check SQLite is available
	if (QSqlDatabase::isDriverAvailable("SQLITE"))
	{
		qCritical("error: sqlite is not available");
		return;
	}
	// Create database as SQLite database
	database = QSqlDatabase::addDatabase("QSQLITE");
	// Check if file exists in path
	if (QFile(path).exists())
	{
		// Just open database
		database.setDatabaseName(path);
		return;
	}
}

DataContext::~DataContext()
{
	// Close database when removing project
	database.close();
}

bool DataContext::create()
{
	QSqlQuery query(database);
	return query.exec("create table Info ("
					  "id integer primary key default 0,"
					  "version integer,"
					  "name text,"
					  "create integer default current_timestamp"
					  ")");
}

bool DataContext::isOpen()
{
	return database.isOpen();
}
