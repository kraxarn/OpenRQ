# OpenRQ
### Requirement management application for Windows, Linux and macOS
### [Next version currently in development (v0.2)](https://github.com/kraxarn/OpenRQ/projects/3)

## Rules for Comitting
* Make sure you use the provided Visual Studio Code settings **before comitting any code**.
* **Never** force push, your changes, and possibly others, **will be deleted**.
* Only makes changes to the files **you have created**.
* **Do not** push your code if it **does not compile**.
* **Always** make sure you have pulled the latest changes **before pushing**.
* File names should always match the name of the struct and **only one struct per file**.
* **No third-party libraries are allowed**, with the exception of any Qt 5.12 and standard Go libraries.

## Code Style
* Structs uses pascal case (`StructName`).
* **All** variables, methods and functions use dromedary case (`variableName`).
* File names should be the **same as the struct name**, but in **all lower case** (`structname.go`)
* Keep everything under package `main`.
* Never leave more than one lines empty.
* `//` comment for **every method and function** (with name) as documentation.

## Directories
* `.vscode`: VS Code settings, **do not modify**.
* `/`: Go source files.

## Branches
There are two main branches
* `golang`: The main development branch for most developers.
* `stable`: Where all changes gets merged to, **do not modify or push to**.

## How to Download
Do not use `git clone ...`, use `go get -v github.com/kraxarn/OpenRQ` instead.

Then, run the application using `go run github.com/kraxarn/OpenRQ`.

The downloaded files are (probably) in `~/go/src/github.com/kraxarn/OpenRQ`.

If you want to run/build with proper version information, run this command in the directory of the source:
```
go run -ldflags "-X main.versionTagName=`git describe --abbrev=0 --tags` -X main.versionCommitHash=`git rev-list origin/golang -1`" .
```

## How to Push
* First, pull the latest changes:
* `git pull`
* Add your modified files:
* `git add <file>`
* Add your comment of what you have done:
* `git commit -m "<comment>"`
* Make sure, again, that you have the latest version:
* `git pull`
* Push your changes, if you get a merge conflict, you broke the rules:
* `git push`
* Make sure your changes were pushed to the **golang** branch.