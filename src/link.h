#pragma once

#include <QString>
#include <QColor>

struct Link
{
	/// Name of link to parent/requirement
	QString name;

	/// Color of the link to parent, used when validating
	QRgb color; ///(Black, Red)
};