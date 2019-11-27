#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQuickStyle>

int main(int argc, char *argv[])
{
	// Enable high-dpi scaling
	QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

	// Create main application
	QGuiApplication app(argc, argv);

	// Set QtQuick style
	QQuickStyle::setStyle("Material");

	// Load main engine
	QQmlApplicationEngine engine;
	const QUrl url(QStringLiteral("qrc:/qml/main.qml"));
	QObject::connect(&engine, &QQmlApplicationEngine::objectCreated,
					 &app, [url](QObject *obj, const QUrl &objUrl) {
		if (!obj && url == objUrl)
			QCoreApplication::exit(-1);
	}, Qt::QueuedConnection);
	engine.load(url);

	// Run while app is running
	return app.exec();
}
