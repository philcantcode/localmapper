package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capabilities/local"
	"github.com/philcantcode/localmapper/propositions"
	"github.com/philcantcode/localmapper/utils"
)

func InitServer() {
	utils.Log("Hosting Server at http://localhost:"+utils.Configs["SERVER_PORT"], true)

	router := mux.NewRouter()

	router.HandleFunc("/capability/run", runCapability)
	router.HandleFunc("/capability/get", getCapabilities)
	router.HandleFunc("/capability/update", updateCapability)

	router.HandleFunc("/local/get-network-adapters", local.HTTP_JSON_GetNetworkAdapters)
	router.HandleFunc("/local/get-os-info", local.HTTP_JSON_GetOSInfo)
	router.HandleFunc("/local/get-date-time", local.HTTP_JSON_GetDateTime)
	router.HandleFunc("/local/get-logs", local.HTTP_JSON_GetLogs)
	router.HandleFunc("/local/get-default-ip-gateway", local.HTTP_JSON_GetDefaultGatewayIP)

	router.HandleFunc("/propositions/get-all", propositions.HTTP_JSON_GetPropositions)
	router.HandleFunc("/propositions/accept-defaults", propositions.ProcessAcceptDefaults)

	fileServer := http.FileServer(http.Dir("/"))

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+utils.Configs["SERVER_PORT"], router)
}
