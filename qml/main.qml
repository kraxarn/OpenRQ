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
			id: content
			anchors.fill: parent
			// Create a new card
			function createCard(properties, textContent)
			{
				const componentString =
					("import QtQuick.Controls 2.12\n" +
					 "import QtQuick.Controls.Material 2.12\n" +
					 "Pane {\n" +
					 "	width: 200\n" +
					 "	height: 100\n" +
					 "	%2\n" +
					 "	Material.elevation: 6\n" +
					 "	Label { text: \"%3\"\n }\n }"
					).arg(properties).arg(textContent)

				return Qt.createQmlObject(componentString, content, "dynamicCard")
			}
			function createLink(from, to, offset)
			{
				const componentString =
					"import QtQuick 2.12\n" +
					"import QtQuick.Shapes 1.12\n" +
					"Shape {\n" +
					"	opacity: 0.5\n" +
					"	ShapePath {\n" +
					"		strokeColor: \"#424242\"\n" +
					"		strokeWidth: 2\n" +
					"		startX: (%1 + %2 / 2) + %3\n".arg(from.x).arg(from.width).arg(offset) +
					"		startY: %1 + %2\n".arg(from.y).arg(from.height) +
					"		PathLine {\n" +
					"			x: %1 + %2 / 2\n".arg(to.x).arg(to.width) +
					"			y: %1\n".arg(to.y) +
					"		}\n" +
					"	}\n" +
					"}"
				return Qt.createQmlObject(componentString, content, "dynamicShape")
			}

			// Create all cards and links
			function createItems()
			{
				// Requirement card
				const requirement = createCard(
					"anchors.top: parent.top; anchors.topMargin: 32; anchors.horizontalCenter: parent.horizontalCenter",
					"Requirements should\nrequire a requirement"
				)

				// Left solution card
				const solution1 = createCard(
					"anchors.topMargin: 64; x: parent.width / 2 - width - 64",
					"Solutions should solve\nthe requirements"
				)

				solution1.anchors.top = requirement.bottom

				// Right solution card
				const solution2 = createCard(
					"anchors.top: cardRequirement.bottom; anchors.topMargin: 64; x: parent.width / 2 + 64",
					"Requirements can be solved\nin different ways"
				)
				solution2.anchors.top = requirement.bottom

				// Create links from requirement to solutions
				createLink(requirement, solution1, -8)
				createLink(requirement, solution2, 8)
			}
			// Start creating cards when ready
			Component.onCompleted: createItems()
		}
	}
}
