package main

// Variables set from linker flags
var versionTagName, versionCommitHash string

func main() {
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
