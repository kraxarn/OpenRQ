import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12

ApplicationWindow {
	visible: true
	width: 1024
	height: 768
	title: qsTr("OpenRQ")
	Material.theme: Material.System

	Label {
		anchors.verticalCenter: parent.verticalCenter
		anchors.horizontalCenter: parent.horizontalCenter
		font.pointSize: 20
        text: "please dont touch me, yes i will"
	}
}
