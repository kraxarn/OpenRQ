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
		QString rationale;

		/// Requirement fit criterion
		QString fitCriterion;

		/// Item::saveChanges()
		bool saveChanges();

		/// Item::getChildren()
		QVector<Item*> getChildren();

		/// Item::getParent()
		Item &getParent();

		/// Item::setParent(Item)
		bool setParent(Item &item);

		/// Item::getHash()
		QByteArray getHash();
	};
}
