package handlers

import (
	"text/template"

	"github.com/philcantcode/localmapper/utils"
)

var TemplateLoader *template.Template

func Reload() { // When done, remove calls to reload
	var err error

	// Parse all .gohtml template files
	TemplateLoader, err = template.ParseGlob("application/web/src/*/*.gohtml")
	utils.ErrorFatal("Couldn't load templates", err)
}
