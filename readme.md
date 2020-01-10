# OpenRQ
OpenRQ is a requirement management application that (attempts to) support versioning, branching and hierarchies. 
It runs on Linux, macOS and Windows and is written in Go with the help of Qt and qt (the Go bindings for Qt and Qt itself have the same name).

# Compile and Run
First, install [Go](https://golang.org/), [Qt](https://qt.io) and [qt](https://github.com/therecipe/qt). 
Any Qt version 5.12 or newer should work, just make sure to generate bindings for Qt 5.12.0. 
If you are on Linux, install the packages through your package manager and on macOS, install the 
packages with [Homebrew](https://brew.sh/). 
After that, you can simply download OpenRQ with `go get` and run it with `go run`.

# Repository Structure
The main source code is located in `/` as Go files. Here is what all the files contain:

| File | Content |
| ---- | ------- |
| datacontext.go | Connection between Go code and the underlying project database |
| edititem.go | Dialogs and handling of editing items |
| icons.go | Mapping to system icons and custom bitmap icons |
| main.go | Entry point |
| mainview.go | The main view for drawing the tree |
| mainwindow.go | Everything related to the window, except for the main view |
| project.go | Handling of loading and parsing projects |
| requirement.go | Requirement item |
| settings.go | Application settings |
| solution.go | Solution item |
| tabledata.go | All tables used in the underlying project database |
| updater.go | Update checker, and previously auto updater |
| validation.go | Validation engine |

The other files and folders are as follows:

| File/Folder | Content |
| ----------- | ------- |
| .vscode | Specific settings for Visual Studio Code, the environment used in development |
| .gitignore | Git ignore file for project and temporary files |
| license | GPL 3.0 license |
| readme.md | This file |
| third_party.md | Third-party licenses and links to source code |