let items = []

function createCard(properties, textContent, color)
{
	const componentString =
			("import QtQuick 2.12\n" +
			"import QtQuick.Controls 2.12\n" +
			"import QtQuick.Controls.Material 2.12\n" +
			"Pane {\n" +
			"	width: 200\n" +
			"	height: 100\n" +
			"	%2\n" +
			"	Material.elevation: 6\n" +
			"	background: Rectangle {\n" +
			"		color: \"transparent\"\n" +
			"		border.color: Material.color(Material.%3)\n" +
			"		border.width: 4\n" +
			"		radius: 12\n" +
			"}\n" +
			"	Label { text: \"%4\"\n }\n }"
		).arg(properties).arg(color).arg(textContent)
	return Qt.createQmlObject(componentString, content, "dynamicCard")
}

function createLink(from, to, offset)
{
	let color = from.background.border.color === to.background.border.color ? "Red" : "Green"

	const componentString =
		"import QtQuick 2.12\n" +
		"import QtQuick.Shapes 1.12\n" +
		"import QtQuick.Controls.Material 2.12\n" +
		"Shape {\n" +
		"	opacity: 1\n" +
		"	ShapePath {\n" +
		"		strokeColor: Material.color(Material.%1)\n".arg(color) +
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

function createItem(type)
{
	let pos, text, color

	// Get position
	switch (items.length)
	{
		case 0:
			pos = "x: parent.width / 2 - width / 2; y: 32"
			break
		case 1:
			pos = "x: parent.width / 2 - width - 64; y: 32 + height + 64"
			break
		case 2:
			pos = "x: parent.width / 2 + 64; y: 32 + height + 64"
			break
		default:
			console.log("too many items")
			pos = "x: 0; y: 0"
			break
	}

	// Get text
	text = type

	// Get color
	color = type === "requirement" ? "Purple" : "Blue"

	items.push(createCard(pos, text, color))

	// Example link creation
	if (items.length > 1)
	{
		createLink(items[0], items[items.length - 1], items.length > 2 ? 8 : -8)
	}
}