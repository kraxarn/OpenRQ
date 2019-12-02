#include "datacontext.h"

namespace orq
{
	DataContext::DataContext(QString path)
	{
		// Check SQLite is available
		if (QSqlDatabase::isDriverAvailable("SQLITE"))
		{
			qCritical("error: sqlite is not available");
			return;
		}

		// Create database as SQLite database
		database = QSqlDatabase::addDatabase("QSQLITE");

		// Check if file already existed
		bool fileExists = QFile(path).exists();

		// Open file
		database.setDatabaseName(path);
		database.open();

		// Create if it didn't exist
		if (!fileExists)
		{
			QString fileName = QFileInfo(path).fileName();
			if (fileName.contains("."))
				fileName = fileName.left(fileName.lastIndexOf("."));
			
			create(fileName);
		}
	}

	DataContext::~DataContext()
	{
		// Close connection
		database.close();
	}

	void DataContext::close()
	{
		// Remove all database connections
		for (auto &connection : QSqlDatabase::connectionNames())
		{
			qDebug() << "debug: removing database connection:" << connection;
			QSqlDatabase::removeDatabase(connection);
		}
	}

	/// Testing creating database from JSON data
	bool DataContext::createFromJson(QString projectName)
	{
		// Load file and try to parse
		QFile file(":/json/tables");
		if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
		{
			qCritical() << "error: failed to open json file with table information";
			return false;
		}
		auto json = QJsonDocument::fromJson(file.readAll());
		if (json.isNull())
		{
			qCritical() << "error: failed to parse json with table information";
			return false;
		}

		// Prepare query
		QSqlQuery query(database);

		// Temporary list to transfer QJsonArray to
		QStringList list;

		// Loop through all entries in the json file
		auto tables = json.object().find("tables")->toObject();
		for (auto key : tables.keys())
		{
			// Clear from previoud
			list.clear();

			// Transfer to QStringList
			for (auto row : tables[key].toArray())
				list << row.toString();

			// Execute query, if false, return
			if (!query.exec(QString("create table %1 (%2)").arg(key).arg(list.join(", "))))
			{
				qCritical() << "database error: failed to create" << key << "table";
				return false;
			}
		}

		// Insert Info table
		query.prepare("insert into Info (name) values (:name)");
		query.bindValue(":name", projectName);
		if (!query.exec())
		{
			qCritical() << "database error: failed to insert into Info table";
			return false;
		}

		// All good
		return true;
	}

	bool DataContext::isOpen()
	{
		return database.isOpen();
	}

	bool DataContext::updateItem(Item &item, int projectVersion)
	{
		// Get item type
		auto type = typeid(item) == typeid(Requirement) ? TypeRequirement : TypeSolution;

		// Prepare query
		QSqlQuery query(database);

		// Check if already pending version
		query.prepare("select * from ItemVersions where version = :version and item = :item");
		query.bindValue(":version", projectVersion);
		query.bindValue(":item", item.id);
		query.exec();

		// If we already have a pending version, just update values
		if (type == TypeRequirement)
		{
			// TODO: Temporary C style cast
			auto req = dynamic_cast<Requirement&>(item);

			query.prepare("update Requirements "
				"set description = :description, rationale = :rationale, fitCriterion = :fitCriterion");
			query.bindValue(":description", req.description);
			query.bindValue(":rationale", req.rationale);
			query.bindValue(":fitCriterion", req.fitCriterion);

			// Execute
			return query.exec();
		}

		// Check if it already exists
		
		query.prepare("select * from :table where id = :id");
		query.bindValue(":table", type == TypeRequirement ? "Requirements" : "Solutions");
		query.bindValue(":id", item.id);
		query.exec();

		// Prepare new version
		int newId = 0;
		if (type == TypeRequirement)
		{
			// TODO: Temporary C style cast
			auto req = dynamic_cast<Requirement&>(item);

			query.prepare("insert into Requirements "
				"(description, rationale, fitCriterion) "
				"values (:description, :rationale, :fitCriterion)");
			query.bindValue(":description", req.description);
			query.bindValue(":rationale", req.rationale);
			query.bindValue(":fitCriterion", req.fitCriterion);
			query.exec();

			// Fetch id for newly added requirement
			query.exec("select * from Requirements order by id desc");
			query.next();
			auto latest = Requirement(query);
			newId = latest.id;
		}
		else
		{
			// ...
		}

		return false;
    }

    /// Check if a specified uid is already taken
    /// (helper for DataContext::getItemUid()
    bool uidExists(QSqlQuery query, qint64 uid)
    {
        // Prepare union query
        query.prepare(
            "select count(*) from "
            "(select uid from Requirements union "
            "select uid from Solutions) where uid = :uid"
        );
        // Bind uid value
        query.bindValue(":uid", uid);
        // Execute the query
        query.exec();
        // If the count(*) value is above 0, uid already exists
        return query.value(0).toInt() > 0;
    }

    qint64 DataContext::getItemUid()
    {
        // Id to return later
        qint64 id;

        // Prepare the query
        QSqlQuery query(database);

        // Keep generating a new value until it doesn't exist
        do
            id = static_cast<qint64>(QRandomGenerator::global()->generate64());
        while (uidExists(query, id));

        // Return the new unique value
        return id;
    }
}
