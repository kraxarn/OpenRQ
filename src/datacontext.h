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

	~DataContext();

	/// Check if database is open
	bool isOpen();

	/// Get path to current database
	QString getCurrentPath();

private:
	/// Private connection to database
	QSqlDatabase database;

	/// Try to create a new database
	bool create(QString projectName);
};
