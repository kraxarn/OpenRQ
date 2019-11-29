#pragma once

#include <QString>
#include <QVector>

#include "item.h"

namespace orq
{
	class Solution:Item
	{
	public:
		Solution();

		~Solution();

		QString linkName;

		QString linkColor;

		QByteArray getHash();
	};

}
