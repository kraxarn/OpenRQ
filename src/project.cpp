#include "project.h"

namespace orq
{
	// Define data in this scope
	DataContext *Project::data;

	Project::Project(QString path)
	{
		// Append .orq if needed
		if (!path.endsWith(".orq"))
			path += ".orq";

		// Create database
		data = new DataContext(path);
	}

	Project::~Project()
	{
		// Close database when destroying object
		delete data;
		data = nullptr;
		// Close database connections
		DataContext::close();
	}

	QHash<int, QPair<ItemType, int>> Project::getVersions()
	{
		qFatal("error: Project::getVersions() is not implemented");
	}

	DataContext& Project::getData()
	{
		return *data;
	}
}
