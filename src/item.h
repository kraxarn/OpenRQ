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
		virtual ~Item() = default;

		int id;

		bool shown = true;

		/// Shared item description
		QString description;

		/// Save changes to database and incremeent version
		virtual bool saveChanges() = 0;

		/// Get all children, probably opposite type
		virtual QVector<Item*> getChildren() = 0;

		/// Get the parent, probably opposite type
		virtual Item &getParent() = 0;

		/// Set parent, should be opposite type
		virtual bool setParent(Item &item) = 0;

		/// Get the MD5 hash for the item
		virtual QByteArray getHash() = 0;
	};
}
