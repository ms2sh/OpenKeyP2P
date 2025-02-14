package p2p

import (
	"net"
	"regexp"
)

// Prüft, ob eine Eingabe eine IPv4, IPv6 oder eine Domain ist
func IdentifyAddressType(address string) AddressType {
	// Prüft, ob es eine gültige IP (IPv4 oder IPv6) ist
	ip := net.ParseIP(address)
	if ip != nil {
		if ip.To4() != nil {
			return AddressTypeIPv4Address
		}
		return AddressTypeIPv6Address
	}

	// Regulärer Ausdruck für eine gültige Domain (z.B. "example.com")
	domainRegex := regexp.MustCompile(`^(?i:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)

	if domainRegex.MatchString(address) {
		return AddressTypeDomain
	}

	return AddressTypeUnkown
}
