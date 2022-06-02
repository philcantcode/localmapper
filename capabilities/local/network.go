package local

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

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

	return ipInfo
}

/* HTTP_JSON_GetNetworkAdapters returns all network adapters on the server */
func HTTP_JSON_GetNetworkAdapters(w http.ResponseWriter, r *http.Request) {
	networkAdapters := GetNetworkAdapters()

	json.NewEncoder(w).Encode(networkAdapters)
}
