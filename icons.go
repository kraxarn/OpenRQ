package main

import (
	"encoding/base64"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
)

var iconData = map[string]string{
	// Menu icons
	"file-new": 	"UklGRq4AAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSGEAAAABZ6AmAAg2atNhBfCLIiICZpEunhgxwSaybSePWEUBGMATErIPJNBhAQG/Ix1XPwEKIvo/AfLhDpCYbmBpTQD1kTpkc+dSHamDzJ1Lddj0rzL1AxGRd8gUn90YV20zNdrbAFZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	"file-open": 	"UklGRtoAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSI4AAAABcFTb1rK8tRi6jZ0g7iFcRn8BnNnXgg4uUVw7OFz0+34pEBETQJpMH/G3LXAAZ5sP3Ew2MNnAZAMTGIYs+kDvCzzdTZO8UtlDVZG2g3I34JSlvgjp5gBw8gSFyLkCgJJRgLNokunkCYnMqHIFThkaiKBDfqnkpjqEMInaDaEBZOBW7iCy/5c88O0T/9QNVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	"file-save-as": "UklGRsAAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSHQAAAABcBQAjNMEI7GBhfXbbQWcDp5Hv7tOkNEKYPOugcUMYQmIiAmAVZpVTQaMsibTOhhX0xjXchrAuJbTAMa1nAYwruU0gHE9Abldr+9UofeZIBJCqHD48kLoKxHv+Xg8LOVc8/rHUz3Ru+aUBFlQejmu0KhlA1ZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	"file-quit": 	"UklGRqoAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSF4AAAABYBPZtpPXkxykp4agFFS8CmT86OM2Pw8CImIC6E8vj0x3hByybQigTKCRWbjQDbVmRDuB4SLUCQwnIdlwC/HxhWEqMH0wTElA4IFaqcgqK5VJrOeQbUOHy7F76CcBVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	// Edit menu
	"edit-rename":	"UklGRpwAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSE8AAAABV2Cmbdtg3W+UBmWLiMjW4Q3YRJItZeOV8CWgBCkYQQPpP7Kx+ah/Gojo/wRoMc9HB3SoJnVmjxGdgSHFbJ0XnWHoP1udfe3qAOi4346ZAFZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	"edit-reload":	"UklGRuYAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSJkAAAABcFtr25p8uA7iDknY7l8At5I5fAlCR0nrLn/O2/yC9hExAfQ30037dLLrSSJiKoGeAyFv+hgUAjPcmmYkUuk9sIRKH6ssCY0TVNLOLUtiBqUGmiRmUFvAEDFoHBEVaR9etoDxogZaKgNUZSl+y8kMzmMy6mCVE5XXqJOif4J724qGjcETQ78K+dscQqcZIM1kbX482c0E/U0AVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	// About menu
	"about-app": 		"UklGRsoAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSH0AAAABcF3bjlKFJ9MeDi1pF1qSQw/2byHJrEUDETEB8FYVVu2+t1WgBO50kKPD+WYHO/1Q2SGMCfcQWw9qJACIWaHwkB0+qilugXodLdooAGrTtqBOR4MqHTkKdHjoP8gmhcAhmCY8p5IIyE/Ciz8UgDVQgwn8v18029YU3h/eCgBWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"about-qt": 		"UklGRh4BAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSNEAAAABgGPb2rHn/rGT0naVWqWTXjMwaxszsDp2tm2jsm1/t/jyvfmCAUTEBOBftXRWUQMgR6KKRkAuVTVCIeV3Y8PH39CXZBWmgHaMxMfGAwUjklzXCe0fT0YVd9AilmN0T9I77LkQMWXHIhE+JJlofw4AwyLhAbJsC2kJ9VfvIlE+sjRL7qCFXxoZytKNX0kG+HEHzWJTmuj59XxNKc+Quf0kwgI9AP6vpBcwLcTjjlY3u3vyunfimaSkr0CyR7eBysUQ4TmViwEDEcF8AHBW0Rb/KgBWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"about-licenses":	"UklGRjABAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSOMAAAABgGRt2/Hmib2GpBjFznKyhG7Bto2NZNaZZra+ce1+fYoPxQYiYgLwb7OTu/f3OxNZG9+MpFJOe418NT705YLBXP8Daz6TOZ4loUydcdog//6QgDb5ILO6KfYB3k5x0e4B+jmh22Me6CTJNqDAHd09g4D4JoAg73R3DAEX386ACC91WywC7d9agQpXdWMcBjxt5+etHmCUbbqEfMpBm396jeowRpFT5QX7YOgjn4bLoVBl5Jn0mYCUVEoSFo2DG5eXGwONdtDaxarLx8fL1ZidqdVmc7y+acVqszkBZaJ50+wnAQBWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"about-update": 	"UklGRpAAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSEQAAAABJ0CkbRvb2BwfPiKC5aGgkaTmHhDwIl4AEtr5N0UvAiL6PwG4NpHhB8qua5AbQ1cb10MmDUMl3QjZMBn9zObEbli7HVZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	"about-gc": 		"UklGRpYAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSEkAAAABYBPZthOyrugABaAAB79CFyUKcrQB1SVziXVETID0XKzHiIWBIG3EBBMCgEmIhUR3AQQMOQHQMzysCWrK8kEyOlYIVknDb7ovAFZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	// Tools
	"tools-move":	"UklGRrwAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSHAAAAABN6CgbRum+/yZDN+nERGBD6YssIokKw6pRUDu70g4n4QnAAv4d8FjsBDRf4Vt2zZK967OSOO4RVQHrrNkqJfFoV4WquapPyWiMB/xMB9p03WJLN5O6aJkRNQxXCI6yPiP9Hd2InJl6FDk1gbOh5EAVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	"tools-link":	"UklGRqgAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSFwAAAABcFzbttLcMuDXhHTDqwCpSmqC6OhEvkXHETEB+qvb9QAgARyWOjDh9wG9p+XaOUU6u1I9rZgSjflpp5REUFLB5nNxIP+CKTqi5moui3r8kgh6VM17pkRQ5vf7ulZQOCAmAAAA0AIAnQEqGAAYAD6RRJ1KpaOioagIALASCWkAAD2joAD++M8awAA=",
	// Other menus
	"menu-edit": 	"UklGRqYAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSFkAAAABYBvZtpJHc04T0IZ38HN6/DHOwT+aR4Tatm0YuWXLH/JZuVpZl3BqqC45BVBm/FhKQJvTjJtrx4TXwOrkUmCMLj0CmuBPse8mFBP71I4YOd7XaK9E9vmqBABWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"menu-delete":	"UklGRowAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSD8AAAABT6CgjRTmB1rwgX8byFg2IiLw/T8VSSiKJDUCBywueGIGJXEQ6ff9juj/BNC2AaSaoMpjg6j+hBZesyLXNgQAVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	// Text formatting
	"format-bold": 			"UklGRrYAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSGkAAAABcFpt27I86W9/tq71H+RnBMZ4I+6yhG5AJrOCk0k00vfgL54jYgLwK3nzOM+eRJr8E5rN5sBwb+kAYEimT0t0IjIynEF3c+c/i/0nNJvNzoLM6QAgIEdPm+pEpLIkC7rbG/s5h4mDfxUAVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
	"format-italic": 		"UklGRpQAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSEcAAAABYNzato3D/FGFaEX0IloRvYhWRBW8OXv0kv8bR8QEqKA405N5sZ18rrDJZ3Pxj15mOOT1hMnL8HG3XnZY5NM82N4DzvRqEgBWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"format-underline":		"UklGRqAAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSFQAAAABcFTbdlO9UYgHLAYP2Oigpo2iANDAioBD7/+PI2IC5GMh3ArhP3NEWwnjsZZOr3RPfSyAQovoEswxz0IfxwNYdUx8y6b15awytXOVUfJuTj7hEwFWUDggJgAAANACAJ0BKhgAGAA+kUSdSqWjoqGoCACwEglpAAA9o6AA/vjPGsAA",
	"format-strikethrough":	"UklGRnQAAABXRUJQVlA4WAoAAAAQAAAAFwAAFwAAQUxQSCcAAAABDzD/ERFCTSQpzHfIwAj+dSCDLgcLEf2fAJzVaYJgGc3ZnDHL9gUAVlA4ICYAAADQAgCdASoYABgAPpFEnUqlo6KhqAgAsBIJaQAAPaOgAP74zxrAAA==",
}

var iconNames = map[string]string{
	// File menu
	"file-new": 	"document-new",
	"file-open": 	"document-open",
	"file-save-as":	"document-save-as",
	"file-quit": 	"exit",
	// Edit menu
	"edit-rename":	"text-field",
	"edit-reload":	"reload",
	// About menu
	"about-app": 		"help-about",
	"about-qt": 		"qt",
	"about-licenses":	"license",
	"about-update": 	"download",
	"about-gc": 		"run-clean",
	// Tools
	"tools-move":	"object-move-symbolic",
	"tools-link":	"draw-line",
	// Other menus
	"menu-edit": 	"document-edit",
	"menu-delete":	"delete",
	// Text formatting
	"format-bold": 			"format-text-bold",
	"format-italic": 		"format-text-italic",
	"format-underline": 	"format-text-underline",
	"format-strikethrough": "format-text-strikethrough",
}

func GetIcon(name string) *gui.QIcon {
	// Try to get icon from theme, otherwise, use bitmap fallback
	return gui.QIcon_FromTheme2(iconNames[name], GetBitmapIcon(name))
}

func GetBitmapIcon(name string) *gui.QIcon {
	bitmap := gui.NewQPixmap()
	data, err := base64.StdEncoding.DecodeString(iconData[name])
	if err != nil {
		fmt.Printf("warning: failed to load icon %v: %v", name, err)
	}
	bitmap.LoadFromData(data, uint(len(data)), "webp", core.Qt__AutoColor)
	return gui.NewQIcon2(bitmap)
}