package main

import "github.com/therecipe/qt/core"

// Variables set from linker flags
var versionTagName, versionCommitHash string

func main() {
	// Setup some application variables
	core.QCoreApplication_SetOrganizationName("kraxarn")
	core.QCoreApplication_SetOrganizationDomain("kraxarn.com")
	core.QCoreApplication_SetApplicationName("OpenRQ")

	// Create window and main app
	app, window := NewMainWindow()

	// Create example project
	NewProject("project.orq")

	// Add menu bar
	CreateLayout(window)
	AddToolBar(window)

	// Show the main window
	window.Show()

	// Main Qt event loop
	app.Exec()
}
