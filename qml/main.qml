import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12

ApplicationWindow {
	id: window
	visible: true
	// Set window size
	width: 1280
	height: 720
	// Center on screen
	x: Screen.width  / 2 - width  / 2
	y: Screen.height / 2 - height / 2
	// Set window title
	title: "OpenRQ"
	// Match theme after system setting
	Material.theme: Material.System

	// Main layout
	Page {
		// Main layout should fill the entire window
		anchors.fill: parent
		// Show a tool bar at the top of the window
		header: ToolBar {
			// Row for each item in the tool bar
			RowLayout {
				// Title or logo
				Label {
					text: "OpenRQ"
					font.pointSize: 16
					padding: 12
				}
				// Dropdown menus
				MenuBar {
					// Hide the default background color
					background: Rectangle {
						color: "transparent"
					}
					// File menu
					Menu {
						title: "File"
						MenuItem {
							text: "New..."
						}
						MenuItem {
							text: "Open..."
						}
						MenuItem {
							text: "Save"
						}
						MenuItem {
							text: "Save As..."
						}
						MenuSeparator {}
						MenuItem {
							text: "Close"
							// Close window when close button is clicked
							onClicked: {
								window.close()
							}
						}
					}
				}
			}
		}

		// Temporary label
		Label {
			// Center vertically and horizontally in parent (Page)
			anchors.verticalCenter: parent.verticalCenter
			anchors.horizontalCenter: parent.horizontalCenter
			// Set larget font
			font.pointSize: 20
			// Align center in multiple lines
			horizontalAlignment: Text.AlignHCenter
			// Set text to show
			text: "Welcome!\nNo project opened"
		}
	}
}
