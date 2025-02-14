package p2p

import (
	"fmt"
	"net"
)

func GetIpFromDomain(domainStr string) (net.IP, error) {
	names, err := net.LookupIP(domainStr)
	if err != nil {
		return nil, fmt.Errorf("GetIpFromDomain: " + err.Error())
	}

	if len(names) < 1 {
		return nil, fmt.Errorf("no ip found")
	}

	return names[0], nil
}
