#pragma once

#include <QString>
#include <QVector>

namespace orq
{
	class Item
	{
    public:
        Item();
        ~Item();
        QString Description;
        bool shown;
        bool saveChanges();
        QVector<Item> getChildren();
        Item getParent();
        bool setParent(Item item);
        bool edit();
        bool compare(Item item);
        //QVector<Label> getLabels();

	private:
        int id;
	};
}
