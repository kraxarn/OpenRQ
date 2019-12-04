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

		/// Item ID
		int id = 0;

		/// Item version
		int version = 0;

		/// If item (and children) are shown, default true
		bool shown = true;

		/// Temporary padding after 1 byte bool
		char padding[3];

		/// Unique identifier
		qint64 uid = 0;

		/// Shared item description
		QString description = QString();

		/// Save changes to database and increment version
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
