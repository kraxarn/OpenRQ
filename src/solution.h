#pragma once

#include <QString>
#include <QVector>

#include "item.h"
#include "link.h"

namespace orq
{
	class Solution : public Item
	{
	public:
		/// Create a new empty solution
		Solution();
		
		/// Link info to parent
		struct Link link;
	};
}