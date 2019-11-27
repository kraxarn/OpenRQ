import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12

ApplicationWindow {
	visible: true
	width: 1024
	height: 768
    title: qsTr("Hello World ")
	Material.theme: Material.System

	Label {
		anchors.verticalCenter: parent.verticalCenter
		anchors.horizontalCenter: parent.horizontalCenter
<<<<<<< HEAD
        font.pointSize: 20
        text: ":) Merge Conflict :( william 2 "

=======
		font.pointSize: 20
		text: ":) Merge Conflict :( william 2 matko "
>>>>>>> e965f62b535e51d0956c9159daad6d80a4c2f604
	}

 CheckBox {
	 id: checkBox
	 x: 424
	 y: 453
	 text: qsTr("Check Box")
 }
}
