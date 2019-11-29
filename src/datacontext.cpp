#include "datacontext.h"

namespace orq
{
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
		// Check if file already existed
		bool fileExists = QFile(path).exists();
		// Open file
		database.setDatabaseName(path);
		// Create file and open it
		database.open();
		// Create if it didn't exist
		if (!fileExists)
		{
			QString fileName = QFileInfo(path).fileName();
			if (fileName.contains("."))
				fileName = fileName.right(fileName.lastIndexOf("."));
			
			create(fileName);
		}
	}

	DataContext::~DataContext()
	{
		// Close database when removing project
		database.close();
	}

	bool DataContext::create(QString projectName)
	{
		QSqlQuery query(database);

		// Create Info table
		if (!query.exec(
			"create table Info ("
			"	id integer primary key,"
			"	version integer default 1,"
			"	name text,"
			"	created integer default current_timestamp"
			")"))
		{
			qCritical() << "database error: failed to create Info table";
			return false;
		}

		// Create Solutions table
		if (!query.exec(
			"create table Solutions ("
			"	id integer primary key,"
			"	parent integer,"
			"	label integer,"
			"	description text,"
			"	link text,"
			"	foreign key(parent) references Requirements(id),"
			"	foreign key(label) references Labels(id)"
			")"))
		{
			qCritical() << "database error: failed to create Solutions table";
			return false;
		}
		
		// Create Projects table
		if (!query.exec(
			"create table Projects ("
			"	id integer primary key,"
			"	name text,"
			"	created integer default current_timestamp"
			")"))
		{
			qCritical() << "database error: failed to create Projects table";
			return false;
		}
		
		// Create ItemVersions table
		if (!query.exec(
			"create table ItemVersions ("
			"	id integer primary key,"
			"   version integer,"
			"	item integer,"
			"	type integer,"
			"	foreign key (version) references Projects(id)"
			")"))
		{
			qCritical() << "database error: failed to create ItemVersions table";
			return false;
		}

		// Create Requirements table
		if (!query.exec(
			"create table Requirements ("
			"	id integer primary key,"
			"	parent integer,"
			"	label integer,"
			"	description text,"
			"	rationale text,"
			"	fitCriterion text,"
			"	foreign key(parent) references Solutions(id),"
			"	foreign key(label) references Labels(id)"
			")"))
		{
			qCritical() << "database error: failed to create Requirements table";
			return false;
		}

		// Create LabelItems LableItem
		if (!query.exec(
			"create table LabelItems ("
			"	id integer primary key,"
			"	label integer,"
			"	item integer,"
			"	type integer,"
			"	foreign key (label) references Labels(id)"
			")"))
		{
			qCritical() << "database error: failed to create LabelItems table";
			return false;
		}

		// Create Media table
		if (!query.exec(
			"create table Media ("
			"	id integer primary key,"
			"	parent int not null,"
			"	format text default 'webp',"
			"	data blob,"
			"	foreign key(parent) references Solutions(id)"
			")"))
		{
			qCritical() << "database error: failed to create Media table";
			return false;
		}
		
		// Create Labels table
		if (!query.exec(
				"create table Labels ("
				"	id integer primary key,"
				"	tag text,"
				"	color integer"
				")"))
		{
			qCritical() << "database error: failed to create Labels table";
			return false;
		}

		// Insert Info table
		query.prepare("insert into Info (name) values (:name)");
		query.bindValue(":name", projectName);
		if (!query.exec())
		{
			qCritical() << "database error: failed to insert into Info table";
			return false;
		}

		// All good
		return true;
	}

	bool DataContext::isOpen()
	{
		return database.isOpen();
	}
}