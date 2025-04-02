package p2p

import (
	"net"
	"regexp"
	"strings"
)

// Prüft, ob eine Eingabe eine IPv4, IPv6, eine Domain oder eine Tor v3-Adresse ist
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

	// Prüft auf eine gültige Tor v3-Adresse (.onion)
	if strings.HasSuffix(address, ".onion") && len(address) == 56+6 {
		// Regulärer Ausdruck für eine gültige Tor v3-Adresse
		torV3Regex := regexp.MustCompile(`^[a-z2-7]{56}\.onion$`)
		if torV3Regex.MatchString(address) {
			return AddressTypeOnionV3
		}
	}

	return AddressTypeUnkown
}
