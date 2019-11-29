#ifndef REQUIREMENT_H
#define REQUIREMENT_H
#pragma once

#include <QString>
#include <QVector>
#include "item.h"

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

#endif // REQUIREMENT_H
