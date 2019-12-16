package main

func main() {
	// Create window and main app
	app, window := NewMainWindow()

	// Add menu bar
	CreateLayout(window)
	AddToolBar(window)

	// make the window visible
	window.Show()

	// Create example project
	NewProject("project.orq")

	// start the main Qt event loop
	app.Exec()
}
