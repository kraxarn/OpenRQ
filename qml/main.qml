import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12
import QtQuick.Shapes 1.12

import "qml/main.js" as MainJs

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
					Menu {
						title: "Add"
						MenuItem {
							text: "Requirement"
							onClicked: {
								MainJs.createItem("requirement")
							}
						}
						MenuItem {
							text: "Solution"
							onClicked: {
								MainJs.createItem("solution")
							}
						}
						MenuSeparator {
						}
						MenuItem {
							text: "Link"
						}
					}
					Menu {
						title: "View"
						MenuItem {
							text: "Validation Engine"
							onClicked: {
								validationDrawer.visible = true
							}
						}
					}
				}
			}
		}

		// Main content
		Item {
			id: content
			anchors.fill: parent
			// Validation engine drawer
			Drawer {
				id: validationDrawer
				width: parent.width / 5
				height: parent.height
				edge: Qt.RightEdge
				
				ToolBar {
					width: parent.width
					
					Label {
						text: "Validation Engine"
						padding: 8
						anchors.horizontalCenter: parent.horizontalCenter
					}
				}
			}
		}
	}
}