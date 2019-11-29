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
		Requirement();

		Requirement(QSqlQuery query);

		~Requirement();

		QString rationale;
		QString fitCriterion;

		QByteArray getHash();
	};
}
