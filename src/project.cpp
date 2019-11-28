#include "project.h"

Project::Project(QString path)
{
	data = new DataContext(path);
}

Project::~Project()
{
	// Close database when destroying object
	delete data;
}
