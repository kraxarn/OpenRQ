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
		bool success = false;
		return success;
	}

    QVector<Item> Item::getChildren()
	{
        auto vector = QVector<Item>();
        return vector;
	}

	Item Item::getParent()
	{
        auto item = Item();
		return item;
	}

	bool Item::setParent(Item item)
	{
		bool success = false;
		return success;
	}
}