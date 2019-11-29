#include "item.h"
#include <QVector>
#include "label.h"

namespace orq
{

	Item::Item()
	{

	}

	Item::~Item()
	{

	}


	bool Item::saveChanges()
	{
		bool success = false;
		return success;
	}

	bool Item::compare(Item item)
	{

		bool success = false;
		return success;
	}

	QVector<Item> Item:: getChildren()
	{
		QVector<Item> *vector = new QVector<Item>();
		return *vector;
	}

	Item Item ::getParent()
	{
		Item item = Item();
		return item;
	}

	bool Item::setParent(Item item)
	{
		bool success = false;
		return success;
	}

	bool Item:: edit()
	{
	bool success = false;
	return success;
	}

	QVector<Label> Item:: getLabels()
	{
		QVector<Label> *label = new QVector<Label>();
		return *label;
	}

}

