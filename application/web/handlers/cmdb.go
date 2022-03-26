package handlers

import (
	"net/http"

	"github.com/philcantcode/localmapper/adapters/definitions"
	"github.com/philcantcode/localmapper/utils"
)

func init() {
	Reload()
}

func CMDBPage(w http.ResponseWriter, r *http.Request) {
	err := TemplateLoader.ExecuteTemplate(w, "cmdb", definitions.HTMLContents{Title: "CMDB Inventory", Description: "CMDB Device Inventory"})
	utils.ErrorFatal("CMDB Page error", err)
}
