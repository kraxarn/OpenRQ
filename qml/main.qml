import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12

ApplicationWindow {
	visible: true
	width: 1024
	height: 768
<<<<<<< HEAD
	color: "#ff557f"
	title: qsTr("Hello World ")
=======
	title: qsTr("OpenRQ")
>>>>>>> 6735678d1569ebe550cca63db222448f0f074fc6
	Material.theme: Material.System

	Label {
		anchors.verticalCenter: parent.verticalCenter
		anchors.horizontalCenter: parent.horizontalCenter
		font.pointSize: 20
<<<<<<< HEAD
		text: ":) Merge Conflict :( william 2 "
=======
        text: "please dont touch me"
>>>>>>> 6735678d1569ebe550cca63db222448f0f074fc6
	}
}
