package handlers

import (
	"net/http"

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
