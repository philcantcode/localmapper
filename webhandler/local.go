package webhandler

import (
	"encoding/json"
	"net/http"

	"github.com/philcantcode/localmapper/local"
)

type LocalHandler struct {
}

var Local = LocalHandler{}

/* HTTP_JSON_GetNetworkAdapters returns all operating system info */
func (lcl *LocalHandler) HTTP_JSON_GetOSInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(local.OSInfo())
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func (lcl *LocalHandler) HTTP_JSON_GetDefaultGatewayIP(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(local.GetDefaultIPGateway())
}

/* HTTP_JSON_GetNetworkAdapters returns all network adapters on the server */
func (lcl *LocalHandler) HTTP_JSON_GetNetworkAdapters(w http.ResponseWriter, r *http.Request) {
	networkAdapters := local.GetNetworkAdapters()
	json.NewEncoder(w).Encode(networkAdapters)
}
