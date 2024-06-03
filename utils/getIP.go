package utils

import (
	"log"
	"net"
)

func GetIPAddress(interfaceName string) string {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Printf("Error getting interface %s: %v", interfaceName, err)
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		log.Printf("Error getting addresses for interface %s: %v", interfaceName, err)
		return ""
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		// 忽略IPv6地址，只返回IPv4地址
		if ip != nil && ip.To4() != nil {
			return ip.String()
		}
	}

	return ""
}
