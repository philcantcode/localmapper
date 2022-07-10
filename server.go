package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/philcantcode/localmapper/local"
	"github.com/philcantcode/localmapper/system"
	"github.com/philcantcode/localmapper/tools/nmap"
	"github.com/philcantcode/localmapper/utils"
	"github.com/philcantcode/localmapper/webhandler"
)

func initServer() {
	system.Log(fmt.Sprintf("Hosting Server at http://%s:%s", local.GetDefaultIPGateway().DefaultIP, system.GetConfig("server-port")), true)

	router := mux.NewRouter()

	// Capability - GET
	router.HandleFunc("/capability/get/all", webhandler.Capability.HTTP_JSON_GetAll)
	router.HandleFunc("/capability/get/{id}", webhandler.Capability.HTTP_JSON_GetByID)
	router.HandleFunc("/capability/get/new", webhandler.Capability.HTTP_JSON_BLANK_Capability)
	router.HandleFunc("/capability/get/new/command", webhandler.Capability.HTTP_JSON_BLANK_Command)
	router.HandleFunc("/capability/get/new/param", webhandler.Capability.HTTP_JSON_BLANK_Param)

	// Capability - UPDATE
	router.HandleFunc("/capability/update", webhandler.Capability.HTTP_JSON_Update)

	// Capability - RUN
	router.HandleFunc("/capability/run", webhandler.Capability.HTTP_NONE_Run)

	// Capability - UTILS
	router.HandleFunc("/capability/utils/restore", webhandler.Capability.HTTP_NONE_Restore)

	// Jobs - GET
	router.HandleFunc("/capability/jobs/get/stats", webhandler.Jobs.HTTP_JSON_GetStats)

	// Compatability - GET
	router.HandleFunc("/compatability/get/capabilities/{entityID}", webhandler.Compatability.HTTP_JSON_GET_Capability_ByEntityID)
	router.HandleFunc("/compatability/run/capabilities/{entityID}/{capabilityID}", webhandler.Compatability.HTTP_NONE_RunCapability)

	// System - GET
	router.HandleFunc("/system/get-logs", webhandler.System.HTTP_JSON_GetLogs)

	// System - UTILS
	router.HandleFunc("/system/execute/action/{id}", webhandler.System.HTTP_JSON_ExecuteAction)
	router.HandleFunc("/system/utils/restore", webhandler.System.HTTP_JSON_Restore)

	// Cookbooks - RUN
	router.HandleFunc("/cookbook/run/{ccbi}/{id}", webhandler.Cookbooks.HTTP_JSON_Run_Cookbook)

	// Cookbooks - GET
	router.HandleFunc("/cookbook/get/all", webhandler.Cookbooks.HTTP_JSON_GetAll)
	router.HandleFunc("/cookbook/get/new", webhandler.Cookbooks.HTTP_JSON_GetNew)

	// Cookbooks - UTILS
	router.HandleFunc("/cookbook/utils/restore", webhandler.Cookbooks.HTTP_JSON_Restore)

	router.HandleFunc("/local/get-network-adapters", local.HTTP_JSON_GetNetworkAdapters)
	router.HandleFunc("/local/get-os-info", local.HTTP_JSON_GetOSInfo)
	router.HandleFunc("/local/get-date-time", utils.HTTP_JSON_GetDateTime)
	router.HandleFunc("/local/get-default-ip-gateway", local.HTTP_JSON_GetDefaultGatewayIP)

	// Propositions - GET
	router.HandleFunc("/propositions/get-all", webhandler.Proposition.HTTP_JSON_GetPropositions)
	router.HandleFunc("/propositions/get-count", webhandler.Proposition.HTTP_INT_GetPropositionCount)

	// Propositions - UTILS
	router.HandleFunc("/propositions/refresh", webhandler.Proposition.HTTP_None_Refresh)

	// Propositions - CMDB - RESOLVE
	router.HandleFunc("/propositions/cmdb/resolve/local-identity", webhandler.CMDB.HTTP_None_ResolveLocalIdentityProposition)
	router.HandleFunc("/propositions/cmdb/resolve/ip-conflict", webhandler.CMDB.HTTP_None_ResolveIPConflict)

	// CMDB - GET
	router.HandleFunc("/cmdb/get/{id}", webhandler.CMDB.HTTP_JSON_GetByID)
	router.HandleFunc("/cmdb/get/stats/{id}", webhandler.CMDB.HTTP_JSON_GetStats)
	router.HandleFunc("/cmdb/pending/get/all", webhandler.CMDB.HTTP_JSON_Pending_GetAll)
	router.HandleFunc("/cmdb/inventory/get/all", webhandler.CMDB.HTTP_JSON_Inventory_GetAll)
	router.HandleFunc("/cmdb/inventory/get/local", webhandler.CMDB.HTTP_JSON_GetLocal)
	router.HandleFunc("/cmdb/inventory/get/type/{type}", webhandler.CMDB.HTTP_JSON_GetByType)

	// CMDB - PUT
	router.HandleFunc("/cmdb/inventory/put/vlan", webhandler.CMDB.HTTP_NONE_NewVLAN)

	// CMDB - Pending (Approve / Deny)
	router.HandleFunc("/cmdb/pending/approve", webhandler.CMDB.HTTP_NONE_ApprovePending)
	router.HandleFunc("/cmdb/pending/deny", webhandler.CMDB.HTTP_NONE_DenyPending)
	router.HandleFunc("/cmdb/pending/deny/all", webhandler.CMDB.HTTP_NONE_DenyAll)

	// CMDB - UTILS
	router.HandleFunc("/cmdb/utils/restore", webhandler.CMDB.HTTP_JSON_Restore)

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
