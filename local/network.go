package local

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/jackpal/gateway"
	"github.com/philcantcode/localmapper/system"
)

type DefaultIPGateway struct {
	DefaultIP      string
	DefaultGateway string
}

type NetworkAdapter struct {
	Label string
	IP    string
	IP6   string
	MAC   string
	MAC6  string
}

func GetNetworkAdapters() []NetworkAdapter {
	var adapters = []NetworkAdapter{}

	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		adapter := NetworkAdapter{Label: i.Name, MAC: i.HardwareAddr.String()}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				adapter.IP = fmt.Sprintf("%s", v.IP)
			case *net.IPAddr:
				adapter.IP6 = fmt.Sprintf("%s", v.IP)
			}

		}

		adapters = append(adapters, adapter)
	}

	system.Log("Getting a list of network adapters & IP addresses.", false)

	return adapters
}

/* GenerateListOfGatewaysFromNetworkAdapters calculates the first IP on every
   network adapter attached. */
func GenerateListOfGatewaysFromNetworkAdapters() []NetworkAdapter {
	netAdapters := GetNetworkAdapters()

	// {Adapter Name : IP Address}
	for key, addr := range netAdapters {
		ip := []byte(net.ParseIP(addr.IP).To4())
		gateway := []byte{ip[0], ip[1], ip[2], 1}

		netAdapters[key].IP = string(net.IPv4(ip[0], ip[1], ip[2], 1).String())
		system.Log(fmt.Sprintf("Calculating gateway: %d -> %d", ip, gateway), false)
	}

	return netAdapters
}

/* HTTP_JSON_GetNetworkAdapters returns all network adapters on the server */
func HTTP_JSON_GetNetworkAdapters(w http.ResponseWriter, r *http.Request) {
	networkAdapters := GetNetworkAdapters()

	json.NewEncoder(w).Encode(networkAdapters)
}

func GetDefaultIPGateway() DefaultIPGateway {
	defaultIP, err := gateway.DiscoverInterface()
	system.Fatal("Could not find the default IP", err)

	defaultGateway, err := gateway.DiscoverGateway()
	system.Fatal("Could not find the default Gateway", err)

	gatewayIP := DefaultIPGateway{DefaultIP: defaultIP.String(), DefaultGateway: defaultGateway.String()}

	return gatewayIP
}

/* HTTP_JSON_GetDefaultGatewayIP both the deafult IP and the Gateway */
func HTTP_JSON_GetDefaultGatewayIP(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(GetDefaultIPGateway())
}
