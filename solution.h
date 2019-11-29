#pragma once

#include <QString>
#include <QVector>
#include "src/item.h"
#include  "media.h"

namespace orq
{
	class Solution:Item
	{
		public:
			Solution();
			~Solution();
			QString linkName;
			QString linkColor;
			QVector<Media> getMedia();
	};

}
