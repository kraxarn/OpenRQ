#pragma once

#include <QGuiApplication>

#include "datacontext.h"
#include "itemtype.h"

namespace orq
{
	class Project
	{
	public:
		/// Open a new project
		Project(QString path);

		/// Close project
		~Project();

		/// Get all existing versions
		QHash<int, QPair<ItemType, int>> getVersions();

		/// Get DataContext instance
		DataContext &getData();

	private:
		/// Project/database has been opened
		bool open;

		/// Data context for data access
		DataContext *data;
	};
}
