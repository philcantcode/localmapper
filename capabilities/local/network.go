package local

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/jackpal/gateway"
	"github.com/philcantcode/localmapper/utils"
)

type DefaultIPGateway struct {
	DefaultIP      string
	DefaultGateway string
}

func GetNetworkAdapters() map[string]string {
	var ipInfo = make(map[string]string)

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ipInfo[i.Name] = fmt.Sprintf("%s", ip)
		}
	}

	utils.Log("Getting a list of network adapters & IP addresses.", false)

	return ipInfo
}

/* HTTP_JSON_GetNetworkAdapters returns all network adapters on the server */
func HTTP_JSON_GetNetworkAdapters(w http.ResponseWriter, r *http.Request) {
	networkAdapters := GetNetworkAdapters()

	json.NewEncoder(w).Encode(networkAdapters)
}

func GetDefaultIPGateway() DefaultIPGateway {
	defaultIP, err := gateway.DiscoverInterface()
	utils.ErrorFatal("Could not find the default IP", err)

	defaultGateway, err := gateway.DiscoverGateway()
	utils.ErrorFatal("Could not find the default Gateway", err)

	gatewayIP := DefaultIPGateway{DefaultIP: defaultIP.String(), DefaultGateway: defaultGateway.String()}

	return gatewayIP
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetDefaultGatewayIP(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetDefaultIPGateway())
}
