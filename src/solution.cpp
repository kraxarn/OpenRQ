#include "solution.h"
#include "item.h"
#include "media.h"

namespace orq
{



	Solution::Solution()
	{

	}

	Solution::~Solution()
	{

	}


	QVector<Media> Solution::getMedia()
	{
		QVector<Media> *media = new QVector<Media>();
		return *media;
	}

}
