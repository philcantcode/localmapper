package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/application/web/handlers"
	"github.com/philcantcode/localmapper/utils"
)

func InitServer() {
	utils.Log("Hosting Server at http://localhost:"+utils.Configs["SERVER_PORT"], true)

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.IndexPage)
	router.HandleFunc("/capability/run", runCapability)
	router.HandleFunc("/capability/get", getCapabilities)
	router.HandleFunc("/capability/update", updateCapability)

	fileServer := http.FileServer(http.Dir("application/web/src"))

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+utils.Configs["SERVER_PORT"], router)
}
