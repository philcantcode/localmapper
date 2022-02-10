package localhost

import (
	"fmt"
	"net"
)

func IpInfo() map[string]string {
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
