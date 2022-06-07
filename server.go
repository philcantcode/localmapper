package main

import (
	"net/http"

	"github.com/gorilla/handlers"
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

	router.HandleFunc("/capability/run/cmdb-compatible/{cmbd_id}/{capability_id}", capability.HTTP_JSON_RunCMDBCompatible)
	router.HandleFunc("/capability/run", capability.HTTP_JSON_Run)
	router.HandleFunc("/capability/get/all", capability.HTTP_JSON_GetAll)
	router.HandleFunc("/capability/get/cmdb-compatible/{id}", capability.HTTP_JSON_GetCMDBCompatible)
	router.HandleFunc("/capability/get/{id}", capability.HTTP_JSON_GetByID)
	router.HandleFunc("/capability/update", capability.HTTP_JSON_Update)

	router.HandleFunc("/local/get-network-adapters", local.HTTP_JSON_GetNetworkAdapters)
	router.HandleFunc("/local/get-os-info", local.HTTP_JSON_GetOSInfo)
	router.HandleFunc("/local/get-date-time", local.HTTP_JSON_GetDateTime)
	router.HandleFunc("/local/get-logs", local.HTTP_JSON_GetLogs)
	router.HandleFunc("/local/get-default-ip-gateway", local.HTTP_JSON_GetDefaultGatewayIP)

	router.HandleFunc("/propositions/get-all", proposition.HTTP_JSON_GetPropositions)
	router.HandleFunc("/propositions/accept-defaults", proposition.HTTP_None_AcceptDefault)

	router.HandleFunc("/cmdb/inventory/get/local", cmdb.HTTP_JSON_GetLocal)
	router.HandleFunc("/cmdb/inventory/get/all", cmdb.HTTP_JSON_Inventory_GetAll)
	router.HandleFunc("/cmdb/inventory/get/type/{type}", cmdb.HTTP_JSON_GetByType)
	router.HandleFunc("/cmdb/inventory/get/{id}", cmdb.HTTP_JSON_GetByID)

	router.HandleFunc("/cmdb/pending/get/all", cmdb.HTTP_JSON_Pending_GetAll)
	router.HandleFunc("/cmdb/pending/put", cmdb.HTTP_INSERT_Pending)
	router.HandleFunc("/cmdb/pending/approve", cmdb.HTTP_Pending_Approve)
	router.HandleFunc("/cmdb/pending/deny/all", cmdb.HTTP_Pending_DenyAll)
	router.HandleFunc("/cmdb/pending/deny", cmdb.HTTP_Pending_Deny)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	fileServer := http.FileServer(http.Dir("/"))

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+utils.Configs["SERVER_PORT"], cors(router))
}
