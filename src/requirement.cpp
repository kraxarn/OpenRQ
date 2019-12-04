#include "requirement.h"

namespace orq
{
	Requirement::Requirement()
	{

	}

	Requirement::Requirement(QSqlQuery query)
	{
		qFatal("error: Requirement(QSQlQuery) is not implemented");
	}

	bool Requirement::saveChanges()
	{
		qFatal("error: Requirement::saveChanges() is not implemented");
	}

	QVector<Item*> Requirement::getChildren()
	{
		qFatal("error: Requirement::getChildren() is not implemented");
	}

	bool Requirement::addChild(Item &item)
	{

	}
	
	bool Requirement::removeChild(Item &child)
	{

	}

	QByteArray Requirement::getHash()
	{
		// Compute hash for combined data
		QCryptographicHash hash(QCryptographicHash::Md5);
		hash.addData(QByteArray().append(QByteArray().setNum(id)));
		hash.addData(QByteArray().append(shown));
		hash.addData(description.toUtf8());
		hash.addData(rationale.toUtf8());
		hash.addData(fitCriterion.toUtf8());
		
		return hash.result();
	}
}
