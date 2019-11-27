#pragma once

#include <QString>
#include <QVector>



namespace orq
{

	class Version{};
	class Item{};


	class Layer
	{
		public:
			Layer();
			~Layer();
			Item root;
			bool saveChanges();
			bool compare(Version version); //argument type? - version
			bool createItem();
			bool deleteItem(); //argument int - ID
			QString tag;
			bool shown;

		private:
			int id;
	};

}

