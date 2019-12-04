#pragma once

#include <QString>
#include <QVector>
#include <QSqlQuery>

#include "item.h"
#include "link.h"

namespace orq
{
	class Solution : public Item
	{
	public:
		/// Create a new empty solution
		Solution();

		/// Create a solution from a SQL query
		Solution(QSqlQuery query);
		
		/// Link info to parent
		struct Link link;
	};
}