import QtQuick 2.12
import QtQuick.Window 2.12
import QtQuick.Controls 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls.Material 2.12
import QtQuick.Shapes 1.12

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

		// Main content
		Item {
			anchors.fill: parent
			// Requirement card
			Pane {
				id: cardRequirement
				width: 200
				height: 100
				anchors.top: parent.top
				anchors.topMargin: 32
				anchors.horizontalCenter: parent.horizontalCenter
				Material.elevation: 6
				Label {
					width: parent.width
					text: "Requirements should\nrequire a requirement"
				}
			}
			// Left solution card
			Pane {
				id: cardSolution1
				width: 200
				height: 100
				anchors.top: cardRequirement.bottom
				anchors.topMargin: 64
				x: parent.width / 2 - width - 64
				Material.elevation: 6
				Label {
					text: "Solution should solve\nthe requirements"
				}
			}
			// Right solution card
			Pane {
				id: cardSolution2
				width: 200
				height: 100
				anchors.top: cardRequirement.bottom
				anchors.topMargin: 64
				x: parent.width / 2 + 64
				Material.elevation: 6
				Label {
					text: "Requirements can be solved\nin different ways"
				}
			}
			// Requirement -> Solution1
			Shape {
				opacity: 0.25
				ShapePath {
					strokeColor: "#9e9e9e"
					strokeWidth: 2
					capStyle: ShapePath.RoundCap
					startX: (cardRequirement.x + cardRequirement.width / 2) - 4
					startY: cardRequirement.y + cardRequirement.height
					PathLine {
						x: cardSolution1.x + cardSolution1.width / 2
						y: cardSolution1.y
					}
				}
			}
			// Requirement -> Solution2
			Shape {
				opacity: 0.25
				ShapePath {
					strokeColor: "#9e9e9e"
					strokeWidth: 2
					capStyle: ShapePath.RoundCap
					startX: (cardRequirement.x + cardRequirement.width / 2) + 4
					startY: cardRequirement.y + cardRequirement.height
					PathLine {
						x: cardSolution2.x + cardSolution1.width / 2
						y: cardSolution2.y
					}
				}
			}
		}
	}
}
