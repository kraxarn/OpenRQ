package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(1280, 720)
	window.SetWindowTitle("OpenRQ")

	// make the window visible
	window.Show()

	p := NewProject("project.orq")
	fmt.Println("database open:", p.data.Database.IsOpen())

	// start the main Qt event loop
	app.Exec()
}
