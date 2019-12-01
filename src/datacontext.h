#pragma once

#include <QGuiApplication>
#include <QSqlDatabase>
#include <QSqlQuery>
#include <QSqlError>
#include <QFileInfo>
#include <QtDebug>

#include "item.h"
#include "requirement.h"
#include "solution.h"
#include "itemtype.h"

namespace orq
{
	class DataContext
	{
	public:
		/// Try to open/create project file
		DataContext(QString path);

		/// Close database and delete any objects
		~DataContext();

		/// Close the database, call after deconstructing
		static void close();

		/// Check if database is open
		bool isOpen();

		/// Get path to current database
		QString getCurrentPath();

		/// Update or create an item
		bool updateItem(Item &item, int projectVersion);
		
	private:
		/// Private connection to database
		QSqlDatabase database;

		/// Creates database, tables and insert stuff into Info
		bool create(QString projectName);
	};
}
