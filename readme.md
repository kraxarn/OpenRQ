# OpenRQ
### Requirement management application for Windows, Linux and macOS
### [Documentation is available here](https://kraxarn.github.io/OpenRQdoc)

## Rules for Comitting
* Import the C++ and QML code styles **before comitting any code**.
* **Never** force push, your changes, and possibly others, **will be deleted**.
* Only makes changes to the files **you have created**.
* **Do not** push your code if it **does not compile**.
* **Always** make sure you have pulled the latest changes **before pushing**.
* File names should always match the name of the class and **only one class per file**.
* **No third-party libraries are allowed**, with the exception of any Qt 5.12 libraries.
* Always make sure the files are in the **correct folders**.
* **Always** use the Qt libraries **instead of** the C++ standard library when possible.

## Code Style
* Classes uses pascal case (`ClassName`).
* **All** Variables and methods use dromedary case (`variableName`).
* File names should be the **same as the class name**, but in **all lower case** (`classname.cpp`)
* Brackets are **always** on a new line in C++, but **never** on a new line in QML.
* C++ source files use the .cpp extension and C++ headers use the .h extension.
* Always use `#pragma once` as include guard.
* Keep everything under namespace `orq`.
* Only use `auto` when the type is already specified elsewhere in the statement.
* No `using namespace` allowed.
* Use matching aliases for your QML resources.
* Never leave more than one lines empty.
* `///` comment for **every method** as documentation.

## Directories
* `codeStyle`: Coding styles for C++ and QML, **do not modify**.
* `qml`: QtQuick layout files, .qml files.
* `src`: C++ source and header files, .cpp/.h files.

## Branches
There are two main branches
* `master`: The main development branch for most developers.
* `stable`: Where all changes gets merged to, **do not modify or push to**.

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
* Make sure your changes were pushed to the **master** branch.