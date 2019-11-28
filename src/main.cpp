#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQuickStyle>

#include "project.h"

int main(int argc, char *argv[])
{
	// Enable high-dpi scaling
	QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

	// Create main application
	QGuiApplication app(argc, argv);

	// Set QtQuick style
	QQuickStyle::setStyle("Material");

	auto p = Project("ORQ");

	// Load main engine
	QQmlApplicationEngine engine;
	const QUrl url(QStringLiteral("qrc:main"));
	QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
					 &app, [url](QObject *obj, const QUrl &objUrl) {
		if (!obj && url == objUrl)
			QCoreApplication::exit(-1);
	}, Qt::QueuedConnection);
	engine.load(url);

	// Run while app is running
	return app.exec();
}
