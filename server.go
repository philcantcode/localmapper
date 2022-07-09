package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/capability"
	"github.com/philcantcode/localmapper/cmdb"
	"github.com/philcantcode/localmapper/cookbook"
	"github.com/philcantcode/localmapper/local"
	"github.com/philcantcode/localmapper/proposition"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/tools/nmap"
	"github.com/philcantcode/localmapper/utils"
)

func initServer() {
	system.Log(fmt.Sprintf("Hosting Server at http://%s:%s", local.GetDefaultIPGateway().DefaultIP, system.GetConfig("server-port")), true)

	router := mux.NewRouter()

	router.HandleFunc("/system/execute/action/{id}", system.HTTP_JSON_ExecuteAction)
	router.HandleFunc("/system/get-logs", system.HTTP_JSON_GetLogs)
	router.HandleFunc("/system/utils/restore", system.HTTP_JSON_Restore)

	router.HandleFunc("/capability/run/cmdb-compatible/{cmbd_id}/{capability_id}", capability.HTTP_JSON_RunCMDBCompatible)
	router.HandleFunc("/capability/run", capability.HTTP_JSON_Run)
	router.HandleFunc("/capability/get/all", capability.HTTP_JSON_GetAll)
	router.HandleFunc("/capability/get/new", capability.HTTP_JSON_GetNew)
	router.HandleFunc("/capability/get/new/command", capability.HTTP_JSON_GetNew_Command)
	router.HandleFunc("/capability/get/new/param", capability.HTTP_JSON_GetNew_Param)
	router.HandleFunc("/capability/get/cmdb-compatible/{id}", capability.HTTP_JSON_GetCMDBCompatible)
	router.HandleFunc("/capability/get/{id}", capability.HTTP_JSON_GetByID)
	router.HandleFunc("/capability/update", capability.HTTP_JSON_Update)
	router.HandleFunc("/capability/utils/restore", capability.HTTP_JSON_Restore)
	router.HandleFunc("/capability/manager/get-tracking", capability.HTTP_JSON_Lifecycle_Manager_List_All)
	router.HandleFunc("/capability/utils/date-time-graph", capability.HTTP_JSON_GetJobsDateTimeGraph)
	router.HandleFunc("/capability/utils/date-job-type-graph", capability.HTTP_JSON_Lifecycle_Manager_JobTypes)
	router.HandleFunc("/capability/jobs/get-stats", capability.HTTP_JSON_GetJobStats)

	router.HandleFunc("/cookbook/run/{ccbi}/{id}", cookbook.HTTP_JSON_Run_Cookbook)
	router.HandleFunc("/cookbook/get/all", cookbook.HTTP_JSON_GetAll)
	router.HandleFunc("/cookbook/get/new", cookbook.HTTP_JSON_GetNew)
	router.HandleFunc("/cookbook/utils/restore", cookbook.HTTP_JSON_Restore)

	router.HandleFunc("/local/get-network-adapters", local.HTTP_JSON_GetNetworkAdapters)
	router.HandleFunc("/local/get-os-info", local.HTTP_JSON_GetOSInfo)
	router.HandleFunc("/local/get-date-time", utils.HTTP_JSON_GetDateTime)
	router.HandleFunc("/local/get-default-ip-gateway", local.HTTP_JSON_GetDefaultGatewayIP)

	router.HandleFunc("/propositions/get-all", proposition.HTTP_JSON_GetPropositions)
	router.HandleFunc("/propositions/refresh", proposition.HTTP_None_Refresh)
	router.HandleFunc("/propositions/get-count", proposition.HTTP_INT_GetPropositionCount)
	router.HandleFunc("/propositions/resolve/local-identity", cmdb.HTTP_None_ResolveLocalIdentityProposition)
	router.HandleFunc("/propositions/resolve/ip-conflict", cmdb.HTTP_None_ResolveIPConflict)

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

	router.HandleFunc("/tools/nmap/select-logs", nmap.HTTP_SELECT_Logs)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	fileServer := http.FileServer(http.Dir("/"))

	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	http.ListenAndServe(":"+system.GetConfig("server-port"), cors(router))
}
