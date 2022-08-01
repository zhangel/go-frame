package utils

import (
	"errors"
	"net"
)

func HostIp() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return "", err
			}
			if ip.IsLoopback() {
				continue
			}
			if ip.To4() == nil {
				continue
			}
			return ip.String(), nil
		}
	}

	return "", errors.New("can't get host ip")
}
