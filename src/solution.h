#pragma once

#include <QString>
#include <QVector>

#include "item.h"

namespace orq
{
	class Solution : public Item
	{
	public:
		/// Create a new empty solution
		Solution();

		/// Name of link to parent/requirement
		QString linkName;

		/// Color of the link to parent, used when validating
		QString linkColor;
	};
}
