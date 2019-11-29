#pragma once

#include <QGuiApplication>
#include <QSqlDatabase>
#include <QSqlQuery>
#include <QFileInfo>
#include <QtDebug>

class DataContext
{
public:
	/// Try to open/create project file
	DataContext(QString path);

	/// Close connection to database
	~DataContext();

	/// Check if database is open
	bool isOpen();

	/// Get path to current database
	QString getCurrentPath();

private:
	/// Private connection to database
	QSqlDatabase database;

	/// Creates database, tables and insert stuff into Info
	bool create(QString projectName);
};
