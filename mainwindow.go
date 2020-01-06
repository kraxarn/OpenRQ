package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// NewMainWindow creates and lays out the main window
func NewMainWindow() (*widgets.QApplication, *widgets.QMainWindow) {
	// Create basic Qt application
	app := widgets.NewQApplication(len(os.Args), os.Args)
	// Create window
	window := widgets.NewQMainWindow(nil, 0)
	// Set minimum and initial size
	window.SetMinimumSize2(960, 540)
	window.Resize2(1280, 720)
	// Center window on screen
	window.SetGeometry(widgets.QStyle_AlignedRect(
		core.Qt__LeftToRight,
		core.Qt__AlignCenter,
		window.Size(),
		gui.QGuiApplication_Screens()[0].AvailableGeometry()))
	// Set a window title
	window.SetWindowTitle("OpenRQ")
	// Return app and window for use in the main function
	return app, window
}

// Temporary global pointer to the validation engine window for the hide/show button
var dockValidation *widgets.QDockWidget

func AddMenuBar(window *widgets.QMainWindow) {
	// Main menu bar
	menuBar := window.MenuBar()

	// File options
	fileMenu := widgets.NewQMenu2("File", nil)
	// New
	fileNew := fileMenu.AddAction2(gui.QIcon_FromTheme("document-new"), "New...")
	fileNew.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__New))
	fileNew.ConnectTriggered(func(checked bool) {
		fileName := widgets.QFileDialog_GetSaveFileName(window, "New Project",
			core.QStandardPaths_Locate(core.QStandardPaths__DocumentsLocation, "", 1),
			"OpenRQ Project(*.orq)", "", 0)
		if len(fileName) > 0 {
			NewProject(fileName)
			ReloadProject(window)
		}
	})
	// Open
	fileOpen := fileMenu.AddAction2(gui.QIcon_FromTheme("document-open"), "Open...")
	fileOpen.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Open))
	fileOpen.ConnectTriggered(func(checked bool) {
		fileName := widgets.QFileDialog_GetOpenFileName(window, "Open Project",
			core.QStandardPaths_Locate(core.QStandardPaths__DocumentsLocation, "", 1),
			"OpenRQ Project(*.orq *orqz)", "", 0)
		if len(fileName) > 0 {
			if strings.HasSuffix(fileName, ".orqz") {
				result := widgets.QMessageBox_Question(window, "Compressed Project",
					"The project you are trying to load is compressed. " +
						"In order to load the project, it first needs to be decompressed. " +
						"This will replace the compressed project with a decompressed project.\n" +
						"Are you sure you want to continue?",
					widgets.QMessageBox__Yes | widgets.QMessageBox__No, widgets.QMessageBox__Yes)
				if result == widgets.QMessageBox__No {
					return
				}
				// Load compressed project
				if _, err := NewCompressedProject(fileName); err != nil {
					widgets.QMessageBox_Critical(window, "Failed to Load Project", err.Error(),
						widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
			} else if strings.HasSuffix(fileName, ".json") {
				if _, err := NewJSONProject(fileName); err != nil {
					widgets.QMessageBox_Critical(window, "Failed to Load Project", err.Error(),
						widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				}
			} else {
				// If not compressed, just load normal project
				NewProject(fileName)
			}
			ReloadProject(window)
		}
	})
	// Save as
	fileSaveAs := fileMenu.AddAction2(gui.QIcon_FromTheme("document-save-as"), "Save As...")
	fileSaveAs.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__SaveAs))
	fileSaveAs.ConnectTriggered(func(checked bool) {
		// Check if project is loaded to save
		if currentProject == nil {
			widgets.QMessageBox_Information(window, "No Project Loaded",
				"No project is current loaded to save", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			return
		}
		// TODO: Set path of current project as default
		fileName := widgets.QFileDialog_GetSaveFileName(window, "Save Project",
			core.QStandardPaths_Locate(core.QStandardPaths__DocumentsLocation, "", 1),
			"OpenRQ Project(*.orq);;OpenRQ Compressed Project(*.orqz);;JavaScript Object Notation(*.json)", "", 0)
		if len(fileName) > 0 {
			if strings.HasSuffix(fileName, ".json") {
				if data, err := json.MarshalIndent(map[string]interface{}{
					"ProjectName": currentProject.Name(),
					"Tree": Roots(),
				}, "", "\t"); err != nil {
					fmt.Println(err)
				} else {
					if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
						fmt.Println(err)
					}
				}
				return
			}
			if err := currentProject.CopyTo(fileName); err != nil {
				widgets.QMessageBox_Critical(window, "Failed to Save Project",
					err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			}
		}
	})
	// Quit
	fileMenu.AddSeparator()
	fileQuit := fileMenu.AddAction2(gui.QIcon_FromTheme("document-new"), "Quit")
	fileQuit.SetShortcut(gui.NewQKeySequence5(gui.QKeySequence__Quit))
	fileQuit.ConnectTriggered(func(checked bool) {
		window.Close()
	})
	menuBar.AddMenu(fileMenu)

	// Edit menu
	editMenu := widgets.NewQMenu2("Edit", nil)
	// Insert menu
	editInsertMenu := widgets.NewQMenu2("Insert", nil)
	editInsertMenu.SetIcon(gui.QIcon_FromTheme("add"))
	editInsertMenu.AddAction2(gui.QIcon_FromTheme("draw-polygon"), "Problem")
	editInsertMenu.AddAction2(gui.QIcon_FromTheme("draw-polygon-star"), "Solution")
	editInsertMenu.AddAction2(gui.QIcon_FromTheme("draw-line"), "Link")
	editMenu.AddMenu(editInsertMenu)
	// Other options
	editMenu.AddAction2(gui.QIcon_FromTheme("text-field"),
		"Rename Project...").ConnectTriggered(func(checked bool) {
		// Check if project is loaded to rename
		if currentProject == nil {
			widgets.QMessageBox_Information(window, "No Project Loaded",
				"No project is current loaded to rename", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			return
		}
		db := currentProject.Data()
		defer db.Close()
		name := widgets.QInputDialog_GetText(window, "Rename Project", "New project name",
			widgets.QLineEdit__Normal, db.ProjectName(), nil, 0, 0)
		if len(name) > 0 {
			db.SetProjectName(name)
			UpdateWindowTitle(window)
		}
	})
	editMenu.AddSeparator()
	editMenu.AddAction2(gui.QIcon_FromTheme("reload"),
		"Reload Project").ConnectTriggered(func(checked bool) {
		ReloadProject(window)
	})
	editMenu.AddAction2(gui.QIcon_FromTheme("view_text"),
		"Dump tree data").ConnectTriggered(func(checked bool) {
		fmt.Println()
		for key, value := range Tree() {
			fmt.Printf("%v: %v\n", key.ToString(), value)
		}
	})
	// Add to main toolbar
	menuBar.AddMenu(editMenu)

	// View
	viewMenu := widgets.NewQMenu2("View", nil)
	validate := viewMenu.AddAction("Validation Engine")
	validate.SetCheckable(true)
	validate.ConnectTriggered(func(checked bool) {
		if checked {
			dockValidation.Show()
		} else {
			dockValidation.Hide()
		}
	})
	menuBar.AddMenu(viewMenu)

	// About
	aboutMenu := widgets.NewQMenu2("About", nil)
	aboutMenu.AddAction2(
		gui.QIcon_FromTheme("help-about"), "About OpenRQ").ConnectTriggered(func(checked bool) {
			// Add app version information
			aboutMessage := "This version was compiled without version information."
			if len(versionTagName) > 0 && len(versionCommitHash) > 0 {
				aboutMessage = fmt.Sprintf("Version\t%v\nCommit\t%v", versionTagName[1:], versionCommitHash)
			}
			// Add useless version and memory information
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			aboutMessage += fmt.Sprintf("\n\nQt %v, Go %v, %v\nMemory usage: %.2f mb (%.2f mb allocated)",
				core.QLibraryInfo_Version().ToString(), runtime.Version()[2:], runtime.GOARCH,
				float64(mem.TotalAlloc)/1000000, float64(mem.Sys)/1000000)
			// Show simple dialog for now
			widgets.QMessageBox_About(window, "About OpenRQ", aboutMessage)
		})
	aboutMenu.AddAction2(gui.QIcon_FromTheme("qt"), "About Qt").ConnectTriggered(func(checked bool) {
		widgets.QMessageBox_AboutQt(window, "About Qt")
	})
	aboutMenu.AddAction2(gui.QIcon_FromTheme("license"), "Licenses").ConnectTriggered(func(checked bool) {
		gui.QDesktopServices_OpenUrl(
			core.NewQUrl3("https://github.com/kraxarn/OpenRQ/blob/golang/third_party.md", 0))
	})
	aboutMenu.AddSeparator()
	aboutMenu.AddAction2(
		gui.QIcon_FromTheme("download"), "Check for updates").ConnectTriggered(func(checked bool) {
			// Check if version was compiled with version information
			if len(versionCommitHash) <= 0 {
				widgets.QMessageBox_About(window, "Updater",
					"This version was compiled without version information,\nupdater is not available")
				return
			}
			// Actually check for updates
			if IsLatestVersion() {
				widgets.QMessageBox_Information(
					window, "Updater", "You are running the latest version",
					widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
				return
			}
			// New update was found
			if widgets.QMessageBox_Question(
				window, "Updater",
				"New update found, do you want to update now?",
				widgets.QMessageBox__Yes|widgets.QMessageBox__No, widgets.QMessageBox__Yes) == widgets.QMessageBox__Yes {
				if err := Update(); err != nil {
					widgets.QMessageBox_Warning(
						window, "Updater", fmt.Sprintf("Failed to update: %v", err),
						widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
				} else {
					widgets.QMessageBox_Information(
						window, "Updater", "Update successful, restart application to apply changes",
						widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
				}
			}
		})
	aboutMenu.AddAction2(gui.QIcon_FromTheme("run-clean"), "Run GC").ConnectTriggered(func(checked bool) {
		// Get memory information
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		// Fetch how much memory we used before
		memBefore := mem.TotalAlloc
		// Run garbage collector
		runtime.GC()
		// Show dialog how much memory we saved
		runtime.ReadMemStats(&mem)
		widgets.QMessageBox_Information(
			window, "Memory Info", fmt.Sprintf("Freed %.2f kb of memory", float64(mem.TotalAlloc-memBefore)/1000.0),
			widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
	})
	menuBar.AddMenu(aboutMenu)
}

// CreateLayout creates the main layout widgets
func CreateLayout(window *widgets.QMainWindow) {
	// Set view as central widget
	linkBtn := widgets.NewQToolButton(nil)
	view := CreateView(window, linkBtn)
	window.SetCentralWidget(view)
	view.Show()
	// Create validation engine dock widget
	dockValidation = widgets.NewQDockWidget("Validation Engine", window, 0)
	dockValidation.SetWidget(CreateValidationEngineLayout())
	// Hide by default
	dockValidation.Hide()
	// Hide close button as that's done from the tool bar
	dockValidation.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetFloatable)
	// Add dock to main window
	window.AddDockWidget(core.Qt__RightDockWidgetArea, dockValidation)

	// Create item type dock widget
	dockItemType := widgets.NewQDockWidget("Tools", window, 0)
	dockItemType.SetWidget(CreateItemTypeCreator(linkBtn))
	// Hide close button as there's no reason to close it
	dockItemType.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetFloatable)
	// Add dock to main window
	window.AddDockWidget(core.Qt__LeftDockWidgetArea, dockItemType)

	// Create item shape dock widget
	dockItemShape := widgets.NewQDockWidget("Shape", window, 0)
	dockItemShape.SetWidget(CreateItemShapeCreator())
	// Hide close button as there's no reason to close it
	dockItemShape.SetFeatures(widgets.QDockWidget__DockWidgetMovable | widgets.QDockWidget__DockWidgetFloatable)
	window.AddDockWidget(core.Qt__LeftDockWidgetArea, dockItemShape)
}

func CreateVBoxWidget(children ...widgets.QWidget_ITF) *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	for _, child := range children {
		layout.AddWidget(child, 1, 0)
	}
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(layout)
	return widget
}

func LayoutToWidget(vBox *widgets.QVBoxLayout) *widgets.QWidget {
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(vBox)
	widget.SetMaximumWidth(200)
	widget.SetMinimumWidth(150)
	return widget
}

func CreateItemTypeCreator(linkBtn *widgets.QToolButton) *widgets.QToolBar {
	layout := widgets.NewQToolBar2(nil)
	// Requirement/solution selection
	moveBtn := widgets.NewQToolButton(nil)
	moveBtn.SetIcon(gui.QIcon_FromTheme("object-move-symbolic"))
	moveBtn.SetText("Move")
	moveBtn.SetCheckable(true)
	moveBtn.SetChecked(true)
	moveBtn.SetToolTip("Move an already created item")
	layout.AddWidget(moveBtn)
	linkBtn.SetIcon(gui.QIcon_FromTheme("draw-line"))
	linkBtn.SetText("Draw line")
	linkBtn.SetCheckable(true)
	linkBtn.SetToolTip("Create a new link between already created items")
	layout.AddWidget(linkBtn)

	moveBtn.ConnectHitButton(func(pos *core.QPoint) bool {
		if moveBtn.IsChecked() {
			return false
		}
		linkBtn.SetChecked(false)
		return true
	})
	linkBtn.ConnectHitButton(func(pos *core.QPoint) bool {
		if linkBtn.IsChecked() {
			return false
		}
		moveBtn.SetChecked(false)
		return true
	})

	return layout
}

func CreateItemShapeCreator() *widgets.QWidget {
	layout := widgets.NewQVBoxLayout()
	shapeList := widgets.NewQListWidget(nil)
	shapeList.SetDragEnabled(true)
	shapeList.AddItem("Square")
	layout.AddWidget(CreateVBoxWidget(shapeList), 0, 0)
	return LayoutToWidget(layout)
}