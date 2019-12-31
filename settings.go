package main

import "github.com/therecipe/qt/core"

type Settings struct {
	settings *core.QSettings
}

func NewSettings() *Settings {
	set := new(Settings)
	set.settings = core.NewQSettings5(nil)
	return set
}

func (set *Settings) LastProject() string {
	return set.settings.Value("lastProject", core.NewQVariant()).ToString()
}

func (set *Settings) SetLastProject(value string) {
	set.settings.SetValue("lastProject", core.NewQVariant1(value))
}