package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/utils"
)

func initServer() {
	utils.Log("Hosting Server at http://localhost:"+utils.Configs["SERVER_PORT"], true)

	router := mux.NewRouter()

	router.HandleFunc("/capability/run", capability.Run)
	router.HandleFunc("/capability/get", capability.Get)
	router.HandleFunc("/capability/update", capability.Update)

	router.HandleFunc("/local/get-network-adapters", local.HTTP_JSON_GetNetworkAdapters)
	router.HandleFunc("/local/get-os-info", local.HTTP_JSON_GetOSInfo)
	router.HandleFunc("/local/get-date-time", local.HTTP_JSON_GetDateTime)
	router.HandleFunc("/local/get-logs", local.HTTP_JSON_GetLogs)
	router.HandleFunc("/local/get-default-ip-gateway", local.HTTP_JSON_GetDefaultGatewayIP)

	router.HandleFunc("/propositions/get-all", proposition.HTTP_JSON_GetPropositions)
	router.HandleFunc("/propositions/accept-defaults", proposition.ProcessAcceptDefaults)

	router.HandleFunc("/cmdb/get-self", cmdb.HTTP_JSON_GetSelf)

	fileServer := http.FileServer(http.Dir("/"))

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+utils.Configs["SERVER_PORT"], router)
}
