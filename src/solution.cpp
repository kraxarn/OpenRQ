#include "solution.h"
#include "item.h"
#include "project.h"

namespace orq
{
	Solution::Solution()
	{

	}

	Solution::Solution(QSqlQuery query)
	{
		qFatal("error: Solution(QSQlQuery) is not implemented");
	}

	QVector<Item*> Solution::getChildren()
	{
		qFatal("error: Solution::getChildren() is not implemented");
	}

	bool Solution::addChild(Item &child)
	{
		return Project::getData().addChild(*this, child);
	}

	bool Solution::removeChild(Item &child)
	{
		return Project::getData().removeChild(*this, child);
	}

	bool Solution::saveChanges()
	{
		qFatal("error: Solution::saveChange() is not implemented");
	}

	QByteArray Solution::getHash()
	{
		// Compute hash for combined data
		QCryptographicHash hash(QCryptographicHash::Md5);
		hash.addData(QByteArray().append(QByteArray().setNum(id)));
		hash.addData(QByteArray().append(shown));
		hash.addData(description.toUtf8());
		
		return hash.result();
	}
}