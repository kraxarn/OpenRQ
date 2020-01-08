package main

import (
	"encoding/base64"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"golang.org/x/tools/go/ssa/interp/testdata/src/runtime"
)

var iconData = map[string]string{
	// Menu icons
	"file-new": "UklGRq4AAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSGEAAAABZ6AmAAg2atNhBfCLIiICZpEunhgxwSaybSePWEUBGMATErIPJNBhAQG/Ix1XPwEKIvo/AfLhDpCYbmBpTQD1kTpkc+dSHamDzJ1Lddj0rzL1AxGRd8gUn90YV20zNdrbAFZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	// Font formatting
	"format-bold": 			"UklGRrYAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSGkAAAABcFpt27I86W9/tq71H+RnBMZ4I+6yhG5AJrOCk0k00vfgL54jYgLwK3nzOM+eRJr8E5rN5sBwb+kAYEimT0t0IjIynEF3c+c/i/0nNJvNzoLM6QAgIEdPm+pEpLIkC7rbG/s5h4mDfxUAVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	"format-italic": 		"UklGRpQAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSEcAAAABYNzato3D/FGFaEX0IloRvYhWRBW8OXv0kv8bR8QEqKA405N5sZ18rrDJZ3Pxj15mOOT1hMnL8HG3XnZY5NM82N4DzvRqEgBWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"format-underline":		"UklGRqAAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSFQAAAABcFTbdlO9UYgHLAYP2Oigpo2iANDAioBD7/+PI2IC5GMh3ArhP3NEWwnjsZZOr3RPfSyAQovoEswxz0IfxwNYdUx8y6b15awytXOVUfJuTj7hEwFWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"format-strikethrough":	"UklGRnQAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSCcAAAABDzD/ERFCTSQpzHfIwAj+dSCDLgcLEf2fAJzVaYJgGc3ZnDHL9gUAVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
}

var iconNames = map[string]string{
	// Menu icons
	"file-new": "document-new",
}

func GetIcon(name string) *gui.QIcon {
	// On Linux, just load icon from theme
	if runtime.GOOS == "linux" {
		//return gui.QIcon_FromTheme(iconNames[name])
	}
	// On other platforms, decode and load image data
	pixmap := gui.NewQPixmap()
	data, err := base64.StdEncoding.DecodeString(iconData[name])
	if err != nil {
		fmt.Printf("warning: failed to load icon %v: %v", name, err)
	}
	pixmap.LoadFromData(data, uint(len(data)), "webp", core.Qt__AutoColor)
	return gui.NewQIcon2(pixmap)
}