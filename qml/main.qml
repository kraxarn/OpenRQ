import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12

ApplicationWindow {
	visible: true
	width: 640
	height: 480
    title: qsTr("Hello World ")
	Material.theme: Material.System

	Label {
		anchors.verticalCenter: parent.verticalCenter
		anchors.horizontalCenter: parent.horizontalCenter
        font.pointSize: 20
        text: ":) Merge Conflict :( william 2 "
	}
}
