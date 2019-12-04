#pragma once

#include <QString>
#include <QVector>
#include <QSqlQuery>

#include "item.h"

namespace orq
{
	class Requirement : public Item
	{
	public:
		/// Create a new empty requirement
		Requirement();

		/// Construct a requirement from a SQL query
		Requirement(QSqlQuery query);

		/// Requirement rationale
		QString rationale = QString();

		/// Requirement fit criterion
		QString fitCriterion = QString();

		/// Item::saveChanges()
		bool saveChanges();

		QVector<Item*> getChildren();

		bool addChild(Item &item);

		bool removeChild(Item &child);

		/// Item::getHash()
		QByteArray getHash();
	};
}
