#include "include/Layer.h"

namespace orq
{

	Layer::Layer()
	{

	}

	Layer::~Layer()
	{

	}


	bool Layer::saveChanges()
	{
		bool success = false;
		return success;
	}

	bool Layer::compare(Version version)
	{ //argument type? - version


		bool success = false;
		return success;
	}

	bool Layer::createItem()
	{
		bool success = false;
		return success;
	}

	bool Layer::deleteItem()
	{ //argument int - ID
		bool success = false;
		return success;
	}

	QString Layer::getTagVersion() const
	{
		return tagVersion;
	}

	void Layer::setTagVersion(const QString &value)
	{
		tagVersion = value;
	}

}
