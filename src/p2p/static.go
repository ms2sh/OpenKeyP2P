package p2p

import "sync"

var (
	nodeConnections map[ConnectionId]*NodeP2PConnection
	controlLock     *sync.Mutex = new(sync.Mutex)
	wasSetuped      bool        = false
)

const (
	AddressTypeIPv4Address AddressType = "ipv4"
	AddressTypeIPv6Address AddressType = "ipv6"
	AddressTypeDomain      AddressType = "domain"
	AddressTypeUnkown      AddressType = "unkown"
)
