#include "requirement.h"

namespace orq
{
	Requirement::Requirement()
	{

	}

	Requirement::~Requirement()
	{

	}

	QByteArray Requirement::getHash()
	{
		// Compute hash for combined data
		QCryptographicHash hash(QCryptographicHash::Md5);
		hash.addData(description.toUtf8());
		hash.addData(QByteArray().append(shown));
		hash.addData(QByteArray().append(id));
		return hash.result();
	}
}
