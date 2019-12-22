package main

// Variables set from linker flags
var versionTagName string
var versionCommitHash string

func main() {
	// Create window and main app
	app, window := NewMainWindow()

	// Add menu bar
	CreateLayout(window)
	AddToolBar(window)

	// Show the main window
	window.Show()

	// Create example project
	NewProject("project.orq")

	// Main Qt event loop
	app.Exec()
}
