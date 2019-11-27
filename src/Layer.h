#ifndef LAYER_H //#ifdef if macro is declared
#define LAYER_H

#include <QString>

class Layer
{
public:
	Layer();
	bool saveChanges();
	bool compare(); //argument type? - version
	bool createItem();
	bool deleteItem(); //argument int - ID
	QString tag;
	bool shown;

private:
	int id;
};

#endif // LAYER_H
