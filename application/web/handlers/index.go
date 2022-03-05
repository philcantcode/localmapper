package handlers

import (
	"net/http"
	"text/template"

	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
)

func init() {
	Reload()
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	Reload()
	err := TemplateLoader.ExecuteTemplate(w, "index", definitions.HTMLContents{Title: "Dashboard", Description: "Main Dashboard"})
	utils.ErrorFatal("Index Page error", err)
}

var TemplateLoader *template.Template

func Reload() { // When done, remove calls to reload
	var err error

	// Parse all .gohtml template files
	TemplateLoader, err = template.ParseGlob("application/web/src/*/*.gohtml")
	utils.ErrorFatal("Couldn't load templates", err)
}
