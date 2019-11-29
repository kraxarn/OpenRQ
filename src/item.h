#pragma once

#include <QString>
#include <QVector>
#include <QByteArray>
#include <QCryptographicHash>

namespace orq
{
	class Item
	{
    public:
		Item();

		int id;
		
		QString description;

		bool shown;

		bool saveChanges();

		QVector<Item> getChildren();

		Item getParent();

		bool setParent(Item item);

		/// Get the MD5 hash for the item
		//QByteArray getHash();
	};
}
