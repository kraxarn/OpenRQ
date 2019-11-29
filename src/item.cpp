#include <QVector>

#include "item.h"

namespace orq
{
	Item::Item()
	{
		// Set default values
		shown = true;
	}

	bool Item::saveChanges()
	{
		return false;
	}

    QVector<Item> Item::getChildren()
	{
		return QVector<Item>();
	}

	Item Item::getParent()
	{
		return Item();
	}

	bool Item::setParent(Item item)
	{
		return false;
	}
}
