#pragma once

#include <QString>
#include <QVector>
#include "src/item.h"

namespace orq
{
	class Requirement: public Item
	{
	public:
			Requirement();
			~Requirement();
			QString rationale;
			QString fitCriterion;
	};

}

