package p2p

import (
	"errors"
	"fmt"
	"net"

	"github.com/quic-go/quic-go"
)

// Gibt die Lokale IP Adresse einer Quic Verbindung aus
func getLocalIPFromConn(conn quic.Connection) string {
	addr := conn.LocalAddr().(*net.UDPAddr)
	if addr.IP.IsUnspecified() {
		ips, err := net.InterfaceAddrs()
		if err == nil {
			for _, ip := range ips {
				if ipnet, ok := ip.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
		return "127.0.0.1"
	}
	return addr.IP.String()
}

// getLocalIPAndPortFromConn gibt die lokale IP-Adresse, den Port sowie den Hostnamen zurück.
func getLocalIPAndHostFromConn(conn quic.Connection) string {
	addr := conn.LocalAddr().(*net.UDPAddr)
	if addr.IP.IsUnspecified() {
		ips, err := net.InterfaceAddrs()
		if err == nil {
			for _, ip := range ips {
				if ipnet, ok := ip.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						// Holen des Hostnamens anhand der IP
						hostname, _ := getHostnameFromIP(ipnet.IP.String())
						// Rückgabe der IP und des Ports als ein String
						if hostname == "" {
							// Falls keine gültige IP gefunden wird, Rückfall auf 127.0.0.1 mit Port
							hostname, _ := getHostnameFromIP("127.0.0.1")
							return fmt.Sprintf("127.0.0.1:%d (%s)", addr.Port, hostname)
						} else {
							return fmt.Sprintf("%s:%d (%s)", ipnet.IP.String(), addr.Port, hostname)
						}
					}
				}
			}
		}
		// Falls keine gültige IP gefunden wird, Rückfall auf 127.0.0.1 mit Port
		hostname, _ := getHostnameFromIP("127.0.0.1")
		return fmt.Sprintf("127.0.0.1:%d (%s)", addr.Port, hostname)
	}
	// Falls eine spezifizierte IP vorhanden ist, IP und Port als String zurückgeben
	hostname, _ := getHostnameFromIP(addr.IP.String())
	return fmt.Sprintf("%s:%d (%s)", addr.IP.String(), addr.Port, hostname)
}

// getRemoteIPAndHostFromConn gibt die Remote-IP-Adresse, den Port sowie den Hostnamen zurück.
func getRemoteIPAndHostFromConn(conn quic.Connection) string {
	addr := conn.RemoteAddr().(*net.UDPAddr)
	hostname, _ := getHostnameFromIP(addr.IP.String())
	return fmt.Sprintf("%s:%d (%s)", addr.IP.String(), addr.Port, hostname)
}

// getHostnameFromIP versucht, den Hostnamen anhand der IP-Adresse zu ermitteln.
func getHostnameFromIP(ip string) (string, error) {
	// Reverse-DNS Lookup, um den Hostnamen zu ermitteln
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err // Wenn kein Hostname gefunden wurde, geben wir einen Fehler zurück
	}
	if len(names) > 0 {
		return names[0], nil // Wir geben den ersten Hostnamen zurück
	}
	return ip, nil // Wenn kein Hostname gefunden wird, zurückgeben von ""
}

// GetInterfaceByIP ermittelt das Netzwerkinterface anhand einer gegebenen IP-Adresse
func getInterfaceByIP(ipAddress string) (*net.Interface, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return nil, errors.New("ungültige IP-Adresse")
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		// Debugging: Alle Adressen des Interfaces ausgeben
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				// IPv4 oder IPv6 Adressen
				if v.IP.Equal(ip) {
					return &iface, nil
				}
			case *net.IPAddr:
				// IPv4 oder IPv6 Adresse
				if v.IP.Equal(ip) {
					return &iface, nil
				}
			}
		}
	}

	// Debugging: IP-Adresse, die nicht gefunden wurde
	fmt.Println("Local IP:", ipAddress)
	return nil, errors.New("kein passendes Interface gefunden")
}
