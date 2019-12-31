package main

import (
	"os"

	"github.com/therecipe/qt/core"
)

type Settings struct {
	settings *core.QSettings
}

func NewSettings() *Settings {
	set := new(Settings)
	set.settings = core.NewQSettings5(nil)
	return set
}

func (set *Settings) LastProject() string {
	project := set.settings.Value("lastProject", core.NewQVariant()).ToString()
	_, err := os.Stat(project)
	if os.IsNotExist(err) {
		return ""
	}
	return project
}

func (set *Settings) SetLastProject(value string) {
	set.settings.SetValue("lastProject", core.NewQVariant1(value))
}