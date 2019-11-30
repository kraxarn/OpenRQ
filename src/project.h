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

	private:
		/// Project/database has been opened
		bool open;

		/// Data context for data access
		DataContext *data;
		
		/// Hash table with <id, <itemType, itemId>>
		QHash<int, QPair<ItemType, int>> *versions;
	};
}
