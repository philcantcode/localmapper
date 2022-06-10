package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/capability/cookbook"
	"github.com/philcantcode/localmapper/capability/local"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/system"
)

func initServer() {
	system.Log("Hosting Server at http://localhost:"+system.Get("server-port"), true)

	router := mux.NewRouter()

	router.HandleFunc("/system/execute/action/{id}", system.HTTP_JSON_ExecuteAction)
	router.HandleFunc("/system/get-logs", system.HTTP_JSON_GetLogs)
	router.HandleFunc("/system/utils/restore", system.HTTP_JSON_Restore)

	router.HandleFunc("/capability/run/cmdb-compatible/{cmbd_id}/{capability_id}", capability.HTTP_JSON_RunCMDBCompatible)
	router.HandleFunc("/capability/run", capability.HTTP_JSON_Run)
	router.HandleFunc("/capability/get/all", capability.HTTP_JSON_GetAll)
	router.HandleFunc("/capability/get/cmdb-compatible/{id}", capability.HTTP_JSON_GetCMDBCompatible)
	router.HandleFunc("/capability/get/{id}", capability.HTTP_JSON_GetByID)
	router.HandleFunc("/capability/update", capability.HTTP_JSON_Update)
	router.HandleFunc("/capability/utils/restore", capability.HTTP_JSON_Restore)

	router.HandleFunc("/cookbook/run/{ccbi}/{id}", cookbook.HTTP_JSON_Run_Cookbook)
	router.HandleFunc("/cookbook/get/all", cookbook.HTTP_JSON_GetAll)
	router.HandleFunc("/cookbook/utils/restore", cookbook.HTTP_JSON_Restore)

	router.HandleFunc("/local/get-network-adapters", local.HTTP_JSON_GetNetworkAdapters)
	router.HandleFunc("/local/get-os-info", local.HTTP_JSON_GetOSInfo)
	router.HandleFunc("/local/get-date-time", local.HTTP_JSON_GetDateTime)
	router.HandleFunc("/local/get-default-ip-gateway", local.HTTP_JSON_GetDefaultGatewayIP)

	router.HandleFunc("/propositions/get-all", proposition.HTTP_JSON_GetPropositions)
	router.HandleFunc("/propositions/accept-defaults", proposition.HTTP_None_AcceptDefault)

	router.HandleFunc("/cmdb/inventory/get/local", cmdb.HTTP_JSON_GetLocal)
	router.HandleFunc("/cmdb/inventory/get/all", cmdb.HTTP_JSON_Inventory_GetAll)
	router.HandleFunc("/cmdb/inventory/get/type/{type}", cmdb.HTTP_JSON_GetByType)

	router.HandleFunc("/cmdb/get/{id}", cmdb.HTTP_JSON_GetByID)

	router.HandleFunc("/cmdb/pending/get/all", cmdb.HTTP_JSON_Pending_GetAll)
	router.HandleFunc("/cmdb/pending/put", cmdb.HTTP_INSERT_Pending_Vlan)
	router.HandleFunc("/cmdb/pending/approve", cmdb.HTTP_Pending_Approve)
	router.HandleFunc("/cmdb/pending/deny/all", cmdb.HTTP_Pending_DenyAll)
	router.HandleFunc("/cmdb/pending/deny", cmdb.HTTP_Pending_Deny)

	router.HandleFunc("/cmdb/identity-confidence/get/{id}", cmdb.HTTP_JSON_IdentityConfidence_Get)
	router.HandleFunc("/cmdb/utils/date-time-graph/get/{id}", cmdb.HTTP_JSON_GetDateTimeGraph)
	router.HandleFunc("/cmdb/utils/restore", cmdb.HTTP_JSON_Restore)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	fileServer := http.FileServer(http.Dir("/"))

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+system.Get("server-port"), cors(router))
}
