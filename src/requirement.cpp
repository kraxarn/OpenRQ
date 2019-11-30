#include "requirement.h"

namespace orq
{
	Requirement::Requirement()
	{

	}

	Requirement::Requirement(QSqlQuery query)
	{

	}

	Requirement::~Requirement()
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
