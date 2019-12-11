package main

import (
	"fmt"
	"os"
)

func main() {
	app, window := NewMainWindow()

	// Create example project
	proj := NewProject("default.orq")
	req := Requirement{}

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
